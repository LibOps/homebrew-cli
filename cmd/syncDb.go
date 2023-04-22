/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/libops/cli/pkg/gcloud"
	"github.com/libops/cli/pkg/libops"
	"github.com/spf13/cobra"
)

// syncDbCmd represents the syncDb command
var syncDbCmd = &cobra.Command{
	Use:   "sync-db",
	Short: "Transfer the database from one environment to another",
	Long: `Transfer the database from one environment to another

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error

		site, _, err := libops.LoadEnvironment(cmd)
		if err != nil {
			log.Fatal(err)
		}

		source, err := cmd.Flags().GetString("source")
		if err != nil {
			return
		}
		target, err := cmd.Flags().GetString("target")
		if err != nil {
			return
		}

		err = libops.PingEnvironment(site, source)
		if err != nil {
			log.Fatal(err)
		}

		err = libops.PingEnvironment(site, target)
		if err != nil {
			log.Fatal(err)
		}

		// get the gcloud id token
		token, err := gcloud.AccessToken()
		if err != nil {
			log.Fatal(err)
		}

		// run the drush command
		exportArgs := []string{
			"sql-dump",
			"-y",
			"--result-file=/tmp/dump.sql",
			"--debug",
		}
		drushCmd := strings.Join(exportArgs, " ")
		libops.IssueCommand(site, source, "drush", drushCmd, token)

		rand := rand.Int()
		fileName := fmt.Sprintf("dump-%s-%d.sql", source, rand)
		gcsObject := fmt.Sprintf("gs://%s-backups/%s", site, fileName)
		uploadArgs := []string{
			"cp",
			"/tmp/dump.sql",
			gcsObject,
		}
		gsutilCmd := strings.Join(uploadArgs, " ")

		libops.IssueCommand(site, source, "gsutil", gsutilCmd, token)

		downloadArgs := []string{
			"cp",
			gcsObject,
			"/tmp/",
		}
		gsutilCmd = strings.Join(downloadArgs, " ")
		libops.IssueCommand(site, target, "gsutil", gsutilCmd, token)

		importArgs := []string{
			"sql-query",
			"-y",
			"--file-delete",
			fmt.Sprintf("--file=/tmp/%s", fileName),
			"--debug",
		}
		drushCmd = strings.Join(importArgs, " ")
		libops.IssueCommand(site, target, "drush", drushCmd, token)
	},
}

func init() {
	rootCmd.AddCommand(syncDbCmd)

	syncDbCmd.Flags().StringP("source", "s", "", "The database that will be exported from")
	syncDbCmd.Flags().StringP("target", "t", "", "The database that will be overwritten")
	syncDbCmd.MarkFlagRequired("source")
	syncDbCmd.MarkFlagRequired("target")
}
