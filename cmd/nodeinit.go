package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// nodeCmd represents the cluster command
var nodeCmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initalize the master node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// Add code here
		// kubeadm init
		fmt.Println("nodeadm init called")
	},
}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().String("token", "", "kubeadm token to be used kubeadm init")
	nodeCmdInit.Flags().String("vip", "192.168.10.5", "VIP ip to be used for multi master setup")
}
