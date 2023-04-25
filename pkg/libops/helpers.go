package libops

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/libops/cli/pkg/gcloud"
	"github.com/spf13/cobra"
)

func LoadEnvironment(cmd *cobra.Command) (string, string, error) {
	site, err := cmd.Flags().GetString("site")
	if err != nil {
		return "", "", err
	}
	env, err := cmd.Flags().GetString("environment")
	if err != nil {
		return site, "", err
	}

	// Perform a DNS lookup on the remote domain to ensure we have sane values
	domain := fmt.Sprintf("%s.remote.%s.libops.site", env, site)
	if _, err := net.LookupHost(domain); err != nil {
		fmt.Println("Error:", err.Error())
		return "", "", fmt.Errorf("Domain %s does not exist. Are site and environment valid?", domain)
	}

	return site, env, nil
}

func IssueCommand(site, env, cmd, args, token string) error {
	var err error
	err = WaitUntilOnline(site, env, token)
	if err != nil {
		return err
	}

	log.Printf("Running `%s %s` on %s %s\n", cmd, args, site, env)
	url := fmt.Sprintf("https://%s.remote.%s.libops.site/%s", env, site, cmd)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(args)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("%s on %s %s returned a non-200: %v", cmd, site, env, resp.StatusCode)
	}
	// print the output to the terminal as it streams in
	for {
		buffer := make([]byte, 1024)
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buffer[:n]))
	}

	return nil
}

func GetToken(cmd *cobra.Command, tokenArg string) (string, error) {
	token, err := cmd.Flags().GetString(tokenArg)
	if err != nil {
		return "", err
	}
	if token == "" {
		token, err = gcloud.AccessToken()
		if err != nil {
			return "", err
		}
	}

	return token, nil
}

func WakeEnvironment(site, env, token string) error {
	url := fmt.Sprintf("https://%s.remote.%s.libops.site/wakeup", env, site)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode > 299 {
		return fmt.Errorf("Environment not able to turn on. %s %s returned a non-200: %v", site, env, resp.StatusCode)
	}

	return nil
}

func WaitUntilOnline(site, env, token string) error {
	var err error
	wakeup := true
	timeout := 3 * time.Minute
	url := fmt.Sprintf("https://%s.remote.%s.libops.site/ping/", env, site)
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	for start := time.Now(); time.Since(start) < timeout; time.Sleep(5 * time.Second) {
		resp, err := client.Do(req)
		if err != nil {
			log.Println(err)
			if wakeup {
				log.Println("Failed to connect, making sure the machine is turned on.")
				err = WakeEnvironment(site, env, token)
				if err != nil {
					log.Println("Trouble turning on machine.")
					log.Println(err)
				} else {
					wakeup = false
				}
			}
			log.Println("Waiting 10 seconds before trying again.")
			time.Sleep(10 * time.Second)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		log.Printf("Received status code %d, retrying...\n", resp.StatusCode)
	}
	log.Println("Timeout exceeded")
	return fmt.Errorf("%s %s not ready after one minute", site, env)
}
