package cmd

import (
	"os"
	"path/filepath"

	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

// nodeCmd represents the cluster command
var nodeCmdReset = &cobra.Command{
	Use:   "reset",
	Short: "Reset node to clean up all kubernetes install and configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var rootDir = filepath.Join(utils.BASE_DIR, utils.KUBERNETES_VERSION)
		kubeadmReset(rootDir)
		cleanup(rootDir)
	},
}

func kubeadmReset(rootDir string) {
	utils.Run(rootDir, "kubeadm", "reset")
}

//TODO needs improvement
func cleanup(rootDir string) {
	os.RemoveAll(rootDir)
	os.RemoveAll(filepath.Join(utils.ETC_DIR, "kubelet.service"))
	os.RemoveAll(filepath.Join(utils.ETC_DIR, "kubelet.service.d"))
	utils.StopAndDisableService("keepalived.service")
	os.RemoveAll(filepath.Join(utils.ETC_DIR, "keepalived.service"))
	os.RemoveAll(filepath.Join(utils.ETC_DIR, "keepalived.service.d"))
	os.RemoveAll("/opt/cni")
}

func init() {
	rootCmd.AddCommand(nodeCmdReset)
}
