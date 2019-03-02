package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/platform9/nodeadm/pkg/logrus"
	executil "github.com/platform9/nodeadm/utils/exec"

	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/systemd"
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
		cleanupDockerImages()
	},
}

func kubeadmReset() {
	log.Infof("[nodeadm:reset] Invoking kubeadm reset")
	cmd := exec.Command(filepath.Join(constants.BaseInstallDir, "kubeadm"), "reset", "--ignore-preflight-errors=all", "--force")
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PATH=%s:%s", constants.BaseInstallDir, os.Getenv("PATH")),
	)
	if err := executil.LogRun(cmd); err != nil {
		log.Warnf("kubeadm reset failed, continuing: %v", err)
	}
}

func cleanupKeepalived() {
	log.Infof("[nodeadm:reset] Stopping & Removing Keepalived")
	if err := systemd.StopIfActive("keepalived.service"); err != nil {
		log.Fatalf("Failed to stop keepalived service: %v", err)
	}
	if err := systemd.DisableIfEnabled("keepalived.service"); err != nil {
		log.Fatalf("Failed to disable keepalived service: %v", err)
	}
	os.RemoveAll(filepath.Join(constants.SystemdDir, "keepalived.service"))
	os.Remove(constants.KeepalivedConfigFilename)
}

func cleanupKubelet() {
	log.Infof("[nodeadm:reset] Stopping & Removing kubelet")
	if err := systemd.StopIfActive("kubelet.service"); err != nil {
		log.Fatalf("Failed to stop kubelet service: %v", err)
	}
	if err := systemd.DisableIfEnabled("kubelet.service"); err != nil {
		log.Fatalf("Failed to disable kubelet service: %v", err)
	}
	failed, err := systemd.Failed("kubelet.service")
	if err != nil {
		log.Fatalf("Failed to check if kubelet service failed: %v", err)
	}
	if failed {
		if err := systemd.ResetFailed("kubelet.service"); err != nil {
			log.Fatalf("Failed to reset failed kubelet service: %v", err)
		}
	}
	os.RemoveAll(filepath.Join(constants.SystemdDir, "kubelet.service"))
	os.RemoveAll(filepath.Join(constants.SystemdDir, "kubelet.service.d"))
}

func cleanupBinaries() {
	log.Infof("[nodeadm:reset] Removing kubernetes binaries")
	os.RemoveAll(filepath.Join(constants.BaseInstallDir, "kubelet"))
	os.RemoveAll(filepath.Join(constants.BaseInstallDir, "kubeadm"))
	os.RemoveAll(filepath.Join(constants.BaseInstallDir, "kubectl"))

	os.RemoveAll(constants.CNIBaseDir)
}

func cleanupNetworking() {
	log.Infof("[nodeadm:reset] Removing flannel state files & resetting networking")
	os.RemoveAll(constants.CNIConfigDir)
	os.RemoveAll(constants.CNIStateDir)
	cmd := exec.Command("ip", "link", "del", "cni0")
	if err := cmd.Run(); err != nil {
		log.Warnf("%q failed, continuing: %v", strings.Join(cmd.Args, " "), err)
	}

	cmd = exec.Command("ip", "link", "del", "flannel.1")
	if err := cmd.Run(); err != nil {
		log.Warnf("%q failed, continuing: %v", strings.Join(cmd.Args, " "), err)
	}
}

func cleanupDockerImages() {
	for _, image := range utils.GetImages() {
		cmd := exec.Command("docker", "rmi", image)
		if err := cmd.Run(); err != nil {
			log.Warnf("%q failed, continuing: %v", strings.Join(cmd.Args, " "), err)
		}
	}
}

func init() {
	rootCmd.AddCommand(nodeCmdReset)
}
