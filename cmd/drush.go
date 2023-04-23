/*
Copyright © 2023 Joe Corall <joe@libops.io>
*/
package cmd

import (
	"log"
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
Info:
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

		// get the gcloud id token
		token, err := gcloud.AccessToken()
		if err != nil {
			log.Fatal(err)
		}

		err = libops.IssueCommand(site, env, "wakeup", "", token)
		if err != nil {
			log.Fatal(err)
		}

		// run the drush command
		drushCmd := strings.Join(args, " ")
		err = libops.IssueCommand(site, env, "drush", drushCmd, token)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(drushCmd)
}
