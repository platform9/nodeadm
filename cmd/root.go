package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var cfgFile string
var LogLevel string

var rootCmd = &cobra.Command{
	Use: "nodeadm",
	Long: `Tool for Kubernetes node management.
This tool lets you initialize, join and reset a node on
your on-premise Kubernetes cluster.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logLevel, err := log.ParseLevel(LogLevel)
		if err != nil {
			log.Fatalf("Could not parse log level %v", logLevel)
		}
		log.SetLevel(logLevel)
		log.SetOutput(os.Stdout)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&LogLevel, "log-level", "l", "info", "set log level for output, permitted values debug, info, warn, error, fatal and panic")
}
