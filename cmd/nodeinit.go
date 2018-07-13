package cmd

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
)

// nodeCmd represents the cluster command
var nodeCmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initalize the master node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config := apis.NodeadmConfiguration{}
		file := ""
		if len(cmd.Flag("cfg").Value.String()) > 0 {
			file = cmd.Flag("cfg").Value.String()
			bytes, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatalf("Failed to read file %s with error %v\n", file, err)
			}
			yaml.Unmarshal(bytes, &config)
		} else {
			tmpFile, err := ioutil.TempFile("", "nodeadm")
			if err != nil {
				log.Fatalf("Failed to create temp file with error %v", err)
			}
			file = tmpFile.Name()
		}
		config.MasterConfiguration.KubernetesVersion = constants.KUBERNETES_VERSION
		config.MasterConfiguration.NoTaintMaster = true
		kubeadm.SetDefaults_MasterConfiguration(&config.MasterConfiguration)
		bytes, err := yaml.Marshal(config.MasterConfiguration)
		if err != nil {
			log.Fatalf("Failed to marshal master config with err %v\n", err)
		}

		err = ioutil.WriteFile(constants.KUBEADM_CONFIG, bytes, constants.READ)
		if err != nil {
			log.Fatalf("Failed to write file %s with error %v\n", file, err)
		}
		utils.InstallMasterComponents(&config)
		kubeadmInit(constants.KUBEADM_CONFIG)
		networkInit(config)
	},
}

func networkInit(config apis.NodeadmConfiguration) {
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
