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

func IssueCommand(site, env, cmd, args, token string) error {
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
