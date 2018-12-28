package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	log "github.com/platform9/nodeadm/pkg/logrus"
	"github.com/platform9/nodeadm/utils"
	executil "github.com/platform9/nodeadm/utils/exec"
	"github.com/spf13/cobra"
)

// nodeCmd represents the cluster command
var nodeCmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initialize the master node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		config := &apis.InitConfiguration{}
		configPath := cmd.Flag("cfg").Value.String()
		if len(configPath) != 0 {
			config, err = utils.InitConfigurationFromFile(configPath)
			if err != nil {
				log.Fatalf("Failed to read configuration from file %q: %v", configPath, err)
			}
		}
		apis.SetInitDefaults(config)
		if err := apis.SetInitDynamicDefaults(config); err != nil {
			log.Fatalf("Failed to set dynamic defaults: %v", err)
		}
		if errors := apis.ValidateInit(config); len(errors) > 0 {
			log.Error("Failed to validate configuration:")
			for i, err := range errors {
				log.Errorf("%v: %v", i, err)
			}
			os.Exit(1)
		}

		masterConfig, err := yaml.Marshal(config.MasterConfiguration)
		if err != nil {
			log.Fatalf("\nFailed to marshal master config with err %v", err)
		}
		err = ioutil.WriteFile(constants.KubeadmConfig, masterConfig, constants.Read)
		if err != nil {
			log.Fatalf("\nFailed to write file %q with error %v", constants.KubeadmConfig, err)
		}

		utils.InstallMasterComponents(config)

		kubeadmInit(constants.KubeadmConfig)

		log.Infoln("Applying workaround for https://github.com/kubernetes/kubeadm/issues/857")
		if err := ensureKubeProxyRespectsHostoverride(); err != nil {
			log.Fatalf("Failed to apply workaround: %v", err)
		}

		networkInit(config)
	},
}

func networkInit(config *apis.InitConfiguration) {
	file := filepath.Join(constants.CacheDir, constants.FlannelDirName, constants.FlannelManifestFilename)
	podSubnetCIDR := config.MasterConfiguration.Networking.PodSubnet
	if len(podSubnetCIDR) == 0 {
		if value, ok := config.MasterConfiguration.ControllerManagerExtraArgs[constants.ControllerManagerClusterCIDRKey]; ok {
			podSubnetCIDR = value
		}
	}
	log.Infof("Pod network %s", podSubnetCIDR)
	manifestStr := utils.Substitute(file, constants.DefaultPodNetwork, podSubnetCIDR)

	cmd := exec.Command(constants.Sysctl, "net.bridge.bridge-nf-call-iptables=1")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
	}

	cmd = exec.Command(filepath.Join(constants.BaseInstallDir, "kubectl"), fmt.Sprintf("--kubeconfig=%s", constants.AdminKubeconfigFile), "apply", "-f", "-")
	reader := strings.NewReader(manifestStr)
	cmd.Stdin = reader
	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
	}

}

func kubeadmInit(config string) {
	cmd := exec.Command(filepath.Join(constants.BaseInstallDir, "kubeadm"), "init", "--ignore-preflight-errors=all", "--config="+config)
	err := executil.LogRun(cmd)
	if err != nil {
		log.Fatalf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
	}
}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().String("cfg", "", "Location of configuration file")
}
