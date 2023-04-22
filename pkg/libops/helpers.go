package libops

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

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

func PingEnvironment(site, env string) error {
	url := fmt.Sprintf("https://%s.drupal.%s.libops.site/", env, site)
	fmt.Println("Ensuring", url, "is online.")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s is not online. Returned %v", url, resp.StatusCode)
	}

	return nil
}

func IssueCommand(site, env, cmd, args, token string) {
	fmt.Printf("Running `%s %s` on %s %s\n", cmd, args, site, env)
	url := fmt.Sprintf("https://%s.remote.%s.libops.site/%s", env, site, cmd)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(args)))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error running request:", err)
		return
	}
	defer resp.Body.Close()

	// print the output to the terminal as it streams in
	for {
		buffer := make([]byte, 1024)
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		fmt.Print(string(buffer[:n]))
	}
}
