package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
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
		err = ioutil.WriteFile(constants.KUBEADM_CONFIG, masterConfig, constants.READ)
		if err != nil {
			log.Fatalf("Failed to write file %q with error %v\n", constants.KUBEADM_CONFIG, err)
		}

		utils.InstallMasterComponents(config)

		kubeadmInit(constants.KUBEADM_CONFIG)

		networkInit(config)
	},
}

func networkInit(config *apis.InitConfiguration) {
	file := filepath.Join(constants.CONF_INSTALL_DIR, "flannel.yaml")
	log.Printf("Pod network %s\n", config.MasterConfiguration.Networking.PodSubnet)
	utils.ReplaceString(file, constants.DEFAULT_POD_NETWORK, config.MasterConfiguration.Networking.PodSubnet)
	utils.Run(constants.BASE_INSTALL_DIR, "sysctl", "net.bridge.bridge-nf-call-iptables=1")
	utils.Run(constants.BASE_INSTALL_DIR, "kubectl", "--kubeconfig="+"/etc/kubernetes/admin.conf", "apply", "-f", file)
}

func kubeadmInit(config string) {
	utils.Run(constants.BASE_INSTALL_DIR, "kubeadm", "init", "--config="+config)
}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().String("cfg", "", "Location of configuration file")
}
