package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/deprecated"
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

// nodeCmd represents the cluster command
var nodeCmdReset = &cobra.Command{
	Use:   "reset",
	Short: "Reset node to clean up all kubernetes install and configuration",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Fail on first error instead of best effort cleanup
		cleanupKeepalived()
		kubeadmReset()
		cleanupKubelet()
		cleanupBinaries()
		cleanupNetworking()
		//cleanupDockerImages()
	},
}

func kubeadmReset() {
	log.Printf("[nodeadm:reset] Invoking kubeadm reset")
	deprecated.RunBestEffort(constants.BaseInstallDir, "kubeadm", "reset")
}

func cleanupKeepalived() {
	log.Printf("[nodeadm:reset] Stopping & Removing Keepalived")
	utils.StopAndDisableService("keepalived.service")
	os.RemoveAll(filepath.Join(constants.SystemdDir, "keepalived.service"))
	os.RemoveAll(filepath.Join(constants.SystemdDir, "keepalived.conf"))
}

func cleanupKubelet() {
	log.Printf("[nodeadm:reset] Stopping & Removing kubelet")
	utils.StopAndDisableService("kubelet.service")
	os.RemoveAll(filepath.Join(constants.SystemdDir, "kubelet.service"))
	os.RemoveAll(filepath.Join(constants.SystemdDir, "kubelet.service.d"))
	err := utils.ResetFailedService("kubelet")
	if err != nil {
		log.Fatalf("Failed to reset failed kubelet service %v\n", err)
	}
}

func cleanupBinaries() {
	log.Printf("[nodeadm:reset] Removing kubernetes binaries")
	os.RemoveAll(filepath.Join(constants.BaseInstallDir, "kubelet"))
	os.RemoveAll(filepath.Join(constants.BaseInstallDir, "kubeadm"))
	os.RemoveAll(filepath.Join(constants.BaseInstallDir, "kubectl"))

	os.RemoveAll(constants.ConfInstallDir)
	os.RemoveAll(constants.CNIBaseDir)
}

func cleanupNetworking() {
	log.Printf("[nodeadm:reset] Removing flannel state files & resetting networking")
	os.RemoveAll(constants.CNIConfigDir)
	os.RemoveAll(constants.CNIStateDir)
	deprecated.RunBestEffort("", "ip", "link", "del", "cni0")
	deprecated.RunBestEffort("", "ip", "link", "del", "flannel.1")
}

func cleanupDockerImages() {
	for _, image := range utils.GetImages() {
		deprecated.RunBestEffort("", "docker", "rmi", "-f", image)
	}
}

func init() {
	rootCmd.AddCommand(nodeCmdReset)
}
