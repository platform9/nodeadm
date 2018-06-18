package cmd

import (
	"path/filepath"

	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

// nodeCmd represents the cluster command
var nodeCmdJoin = &cobra.Command{
	Use:   "join",
	Short: "Initalize the worker node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var rootDir = filepath.Join(utils.BASE_DIR, utils.KUBERNETES_VERSION)
		utils.Install(utils.KUBERNETES_VERSION, utils.CNI_VERSION, rootDir, nil)
		kubeadmJoin(cmd.Flag("token").Value.String(),
			cmd.Flag("master").Value.String(),
			cmd.Flag("cahash").Value.String(), rootDir)
	},
}

func kubeadmJoin(token, master, cahash, rootDir string) {
	utils.Run(rootDir, "kubeadm", "join", "--token", token, master, "--discovery-token-ca-cert-hash", cahash)
}

func init() {
	rootCmd.AddCommand(nodeCmdJoin)
	nodeCmdJoin.Flags().String("token", "", "kubeadm token to be used for kubeadm join")
	nodeCmdJoin.Flags().String("master", "", "masterIP:masterPort for the master to join")
	nodeCmdJoin.Flags().String("cahash", "", "CA hash")
}
