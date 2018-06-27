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
		kubeadmReset()
		cleanup()
	},
}

func kubeadmReset() {
	utils.Run(utils.BASE_DIR, "kubeadm", "reset")
}

//TODO needs improvement
func cleanup() {
	os.RemoveAll(utils.BASE_DIR)
	os.RemoveAll(utils.CNI_BASE_DIR)
	os.RemoveAll(utils.KUBE_DIR)
	os.RemoveAll(utils.CNI_DIR)
	os.RemoveAll(filepath.Join(utils.SYSTEMD_DIR, "kubelet.service"))
	os.RemoveAll(filepath.Join(utils.SYSTEMD_DIR, "kubelet.service.d"))
	utils.StopAndDisableService("keepalived.service")
	os.RemoveAll(filepath.Join(utils.SYSTEMD_DIR, "keepalived.service"))
	os.RemoveAll(filepath.Join(utils.SYSTEMD_DIR, "keepalived.conf"))
	os.RemoveAll(filepath.Join(utils.SYSTEMD_DIR, "keepalived.service.d"))
	os.RemoveAll("/opt/cni")
}

func init() {
	rootCmd.AddCommand(nodeCmdReset)
}
