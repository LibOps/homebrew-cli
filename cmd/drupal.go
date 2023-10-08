/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/libops/cli/pkg/gcloud"
	"github.com/libops/cli/pkg/libops"
	"github.com/spf13/cobra"
)

type Response struct {
	Url string `json:"url"`
}

// getDrupalCmd represents the drupal command
var getDrupalCmd = &cobra.Command{
	Use:   "drupal",
	Short: "View basic information about your LibOps Drupal deployment.",
	Run: func(cmd *cobra.Command, args []string) {
		site, env, err := libops.LoadEnvironment(cmd)
		if err != nil {
			log.Println("Unable to load environment.")
			log.Fatal(err)
		}

		// get the gcloud id token
		token, err := libops.GetToken(cmd, "token")
		if err != nil {
			log.Fatal(err)
		}

		r := Response{}

		url, err := gcloud.GetCloudRunUrl(site, env)
		if err != nil {
			log.Fatal("Unable to retrieve remote URL")
		}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev", url), nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode > 299 {
			log.Fatalf("Unable to get environment info: %v", resp.StatusCode)
		}

		json.NewDecoder(resp.Body).Decode(&r)
		b, err := json.Marshal(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	},
}

func init() {
	getCmd.AddCommand(getDrupalCmd)
	getDrupalCmd.Flags().StringP("token", "t", "", "(optional/machines-only) The gcloud identity token to access your LibOps environment")
}
