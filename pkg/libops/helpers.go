package libops

import (
	"bytes"
	"fmt"
	"io"
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

	return site, env, nil
}

func IssueCommand(site, env, cmd, args, token string) error {
	var err error
	err = WaitUntilOnline(site, env, token)
	if err != nil {
		return err
	}

	fmt.Printf("Running `%s %s` on %s %s\n", cmd, args, site, env)
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
	url := fmt.Sprintf("https://%s.drupal.%s.libops.site/", env, site)
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	var resp *http.Response
	var err error

	wakeup := true

	timeout := 5 * time.Minute
	deadline := time.Now().Add(timeout)
	for {
		select {
		case <-time.After(time.Until(deadline)):
			fmt.Println("Timeout exceeded")
			return fmt.Errorf("%s %s not ready after five minutes", site, env)
		default:
			resp, err = client.Get(url)
			if err != nil {
				if wakeup {
					err = WakeEnvironment(site, env, token)
					if err != nil {
						fmt.Println("Trouble turning on machine.")
						fmt.Println(err)
					} else {
						wakeup = false
					}
				}
				fmt.Println(err)
				fmt.Println("Waiting 10 seconds before trying again.")
				time.Sleep(10 * time.Second)
				continue
			}
			if resp.StatusCode == http.StatusOK {
				fmt.Println(resp.Status)
				return nil
			}
			fmt.Printf("Received status code %d, retrying...\n", resp.StatusCode)
			time.Sleep(5 * time.Second) // wait 5 seconds before retrying
		}
	}
}
