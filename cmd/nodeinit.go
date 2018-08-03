package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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
	Short: "Initalize the master node with given configuration",
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
			log.Println("Failed to validate configuration:")
			for i, err := range errors {
				log.Printf("%v: %v", i, err)
			}
			os.Exit(1)
		}

		masterConfig, err := yaml.Marshal(config.MasterConfiguration)
		if err != nil {
			log.Fatalf("Failed to marshal master config with err %v\n", err)
		}
		err = ioutil.WriteFile(constants.KubeadmConfig, masterConfig, constants.Read)
		if err != nil {
			log.Fatalf("Failed to write file %q with error %v\n", constants.KubeadmConfig, err)
		}

		utils.InstallMasterComponents(config)

		kubeadmInit(constants.KubeadmConfig)

		log.Println("Applying workaround for https://github.com/kubernetes/kubeadm/issues/857")
		if err := workarounds.EnsureKubeProxyRespectsHostoverride(); err != nil {
			log.Fatalf("Failed to apply workaround: %v", err)
		}

		networkInit(config)
	},
}

func networkInit(config *apis.InitConfiguration) {
	file := filepath.Join(constants.ConfInstallDir, constants.FlannelManifestFilename)
	log.Printf("Pod network %s\n", config.MasterConfiguration.Networking.PodSubnet)
	utils.ReplaceString(file, constants.DefaultPodNetwork, config.MasterConfiguration.Networking.PodSubnet)
	deprecated.Run(constants.BaseInstallDir, "sysctl", "net.bridge.bridge-nf-call-iptables=1")
	deprecated.Run(constants.BaseInstallDir, "kubectl", fmt.Sprintf("--kubeconfig=%s", constants.AdminKubeconfigFile), "apply", "-f", file)
}

func kubeadmInit(config string) {
	deprecated.Run(constants.BaseInstallDir, "kubeadm", "init", "--config="+config)
}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().String("cfg", "", "Location of configuration file")
}
