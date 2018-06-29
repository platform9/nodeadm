package cmd

import (
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

var overwriteSymlink bool

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download components",
	Run: func(cmd *cobra.Command, args []string) {
		utils.DownloadArtifacts()
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
