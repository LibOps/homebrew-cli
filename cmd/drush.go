/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/libops/cli/pkg/gcloud"
	"github.com/libops/cli/pkg/libops"
	"github.com/spf13/cobra"
)

// drushCmd represents the drush command
var drushCmd = &cobra.Command{
	Use:   "drush",
	Short: "Run drush commands on your libops environment",
	Long: `
    Run drush commands on your libops environment.

    Currently only non-interactive drush commands are supported.

    If the drush command asks for input the command will fail.

    Examples:
    libops drush -- sql-query -y --file-delete --file=/tmp/dump.sql
    libops drush -- cr
	# enable diff module on the production environment
	libops drush -e production -- en diff
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		site, env, err := libops.LoadEnvironment(cmd)
		if err != nil {
			log.Fatal(err)
		}

		err = libops.PingEnvironment(site, env)
		if err != nil {
			log.Fatal(err)
		}

		// get the gcloud id token
		token, err := gcloud.AccessToken()
		if err != nil {
			log.Fatal(err)
		}

		// run the drush command
		drushCmd := strings.Join(args, " ")
		fmt.Printf("Running `drush %s` on %s %s\n", drushCmd, site, env)
		url := fmt.Sprintf("https://%s.remote.%s.libops.site/drush", env, site)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(drushCmd)))
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

	},
}

func init() {
	rootCmd.AddCommand(drushCmd)
}
