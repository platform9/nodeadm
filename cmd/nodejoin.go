package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nodeCmd represents the cluster command
var nodeCmdJoin = &cobra.Command{
	Use:   "join",
	Short: "Initalize the worker node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// Add code here
		// kubeadm join
		fmt.Println("nodeadm join called")
	},
}

func init() {
	rootCmd.AddCommand(nodeCmdJoin)
	nodeCmdJoin.Flags().String("token", "", "kubeadm token to be used for kubeadm join")
}
