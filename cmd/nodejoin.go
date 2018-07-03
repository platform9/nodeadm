package cmd

import (
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

// nodeCmd represents the cluster command
var nodeCmdJoin = &cobra.Command{
	Use:   "join",
	Short: "Initalize the worker node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		utils.InstallWorkerComponents()
		kubeadmJoin(cmd.Flag("token").Value.String(),
			cmd.Flag("master").Value.String(),
			cmd.Flag("cahash").Value.String())
	},
}

func kubeadmJoin(token, master, cahash string) {
	utils.Run(utils.BASE_INSTALL_DIR, "kubeadm", "join", "--token", token, master, "--discovery-token-ca-cert-hash", cahash)
}

func init() {
	rootCmd.AddCommand(nodeCmdJoin)
	nodeCmdJoin.Flags().String("token", "", "kubeadm token to be used for kubeadm join")
	nodeCmdJoin.Flags().String("master", "", "masterIP:masterPort for the master to join")
	nodeCmdJoin.Flags().String("cahash", "", "CA hash")
}
