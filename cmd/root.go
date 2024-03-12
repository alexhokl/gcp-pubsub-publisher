package cmd

import (
	"github.com/alexhokl/helper/cli"
	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:          "gcp-pubsub",
	Short:        "An application publish a message to Google Cloud Pub/Sub topic",
	SilenceUsage: true,
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	cli.ConfigureViper(cfgFile, "gcp-pubsub-publisher", false, "")
}
