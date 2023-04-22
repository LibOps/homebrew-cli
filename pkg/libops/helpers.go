package libops

import (
	"fmt"
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
