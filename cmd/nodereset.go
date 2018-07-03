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
	utils.RunBestEffort(utils.BASE_INSTALL_DIR, "kubeadm", "reset")
}

//TODO needs improvement
func cleanup() {
	utils.StopAndDisableService("keepalived.service")
	os.RemoveAll(filepath.Join(utils.SYSTEMD_DIR, "keepalived.service"))
	os.RemoveAll(filepath.Join(utils.SYSTEMD_DIR, "keepalived.conf"))

	utils.StopAndDisableService("kubelet.service")
	os.RemoveAll(filepath.Join(utils.SYSTEMD_DIR, "kubelet.service"))
	os.RemoveAll(filepath.Join(utils.SYSTEMD_DIR, "kubelet.service.d"))

	os.RemoveAll(filepath.Join(utils.BASE_INSTALL_DIR, "kubelet"))
	os.RemoveAll(filepath.Join(utils.BASE_INSTALL_DIR, "kubeadm"))
	os.RemoveAll(filepath.Join(utils.BASE_INSTALL_DIR, "kubectl"))

	os.RemoveAll(utils.KUBE_VERSION_INSTALL_DIR)
	os.RemoveAll(utils.CONF_INSTALL_DIR)
	os.RemoveAll(utils.CNI_BASE_DIR)
	os.RemoveAll("/etc/cni")

	utils.RunBestEffort("", "ip", "link", "del", "cni0")
	utils.RunBestEffort("", "ip", "link", "del", "flannel.1")

	for _, image := range utils.GetImages() {
		utils.RunBestEffort("", "docker", "rmi", "-f", image)
	}

}

func init() {
	rootCmd.AddCommand(nodeCmdReset)
}
