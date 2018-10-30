package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/platform9/nodeadm/pkg/logrus"

	"github.com/ghodss/yaml"

	"github.com/platform9/nodeadm/workarounds"

	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/deprecated"
	"github.com/platform9/nodeadm/utils"
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
		if err := workarounds.EnsureKubeProxyRespectsHostoverride(); err != nil {
			log.Fatalf("Failed to apply workaround: %v", err)
		}

		networkInit(config)
	},
}

func networkInit(config *apis.InitConfiguration) {
	file := filepath.Join(constants.CacheDir, constants.FlannelDirName, constants.FlannelManifestFilename)
	log.Infof("\nPod network %s", config.MasterConfiguration.Networking.PodSubnet)
	cmdString := utils.GetContents(file, constants.DefaultPodNetwork, config.MasterConfiguration.Networking.PodSubnet)
	deprecated.Run(constants.BaseInstallDir, "sysctl", "net.bridge.bridge-nf-call-iptables=1")
	deprecated.PipeCmd(constants.BaseInstallDir, cmdString, "kubectl", fmt.Sprintf("--kubeconfig=%s", constants.AdminKubeconfigFile), "apply", "-f", "-")
}

func kubeadmInit(config string) {
	deprecated.Run(constants.BaseInstallDir, "kubeadm", "init", "--ignore-preflight-errors=all", "--config="+config)
}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().String("cfg", "", "Location of configuration file")
}
