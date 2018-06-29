package cmd

import (
	"fmt"
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

var images bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List components to download",
	Run: func(cmd *cobra.Command, args []string) {
		if images {
			images := utils.GetImages()
			for _, image := range images {
				fmt.Println(image)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVar(&images, "images", false, "set to show list of images")
}
