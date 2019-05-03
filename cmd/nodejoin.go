package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/pkg/kubeadm"
	log "github.com/platform9/nodeadm/pkg/logrus"
	executil "github.com/platform9/nodeadm/utils/exec"

	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

var (
	joinCfgPath string
)

// nodeCmd represents the cluster command
var nodeCmdJoin = &cobra.Command{
	Use:   "join",
	Short: "Initalize the node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		config := &apis.JoinConfiguration{}
		if joinCfgPath != "" {
			config, err = utils.JoinConfigurationFromFile(joinCfgPath)
			if err != nil {
				log.Fatalf("Failed to read configuration from file %q: %v", joinCfgPath, err)
			}
		}
		if errors := apis.ValidateJoin(config); len(errors) > 0 {
			log.Error("Failed to validate configuration:")
			for i, err := range errors {
				log.Errorf("%v: %v", i, err)
			}
			os.Exit(1)
		}
		if err := kubeadm.WriteConfiguration(constants.KubeadmConfig, config.NodeConfiguration); err != nil {
			log.Fatalf("Unable to write kubeadm configuration to %s: %s", constants.KubeadmConfig, err)
		}
		utils.InstallNodeComponents()
		kubeadmJoin()
	},
}

func kubeadmJoin() {
	cmd := exec.Command(filepath.Join(constants.BaseInstallDir, "kubeadm"), "join", "--ignore-preflight-errors=all", fmt.Sprintf("--config=%s", constants.KubeadmConfig))
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PATH=%s:%s", constants.BaseInstallDir, os.Getenv("PATH")),
	)
	log.Infof("Running %q", strings.Join(cmd.Args, " "))
	if err := executil.LogRun(cmd); err != nil {
		log.Fatalf("%q failed: %s", strings.Join(cmd.Args, " "), err)
	}
}

func init() {
	rootCmd.AddCommand(nodeCmdJoin)
	nodeCmdJoin.Flags().StringVar(&joinCfgPath, "cfg", "", "Location of configuration file")
}
