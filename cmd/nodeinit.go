package cmd

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
)

// nodeCmd represents the cluster command
var nodeCmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initalize the master node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.Configuration{}
		file := ""
		if cmd.Flag("cfg") != nil {
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
		config.MasterConfiguration.KubernetesVersion = utils.KUBERNETES_VERSION

		kubeadm.SetDefaults_MasterConfiguration(&config.MasterConfiguration)
		bytes, err := yaml.Marshal(config.MasterConfiguration)
		if err != nil {
			log.Fatalf("Failed to marshal master config with err %v\n", err)
		}

		err = ioutil.WriteFile(utils.KUBEADM_CONFIG, bytes, utils.FILE_MODE)
		if err != nil {
			log.Fatalf("Failed to write file %s with error %v\n", file, err)
		}
		utils.InstallMasterComponents(&config)
		kubeadmInit(utils.KUBEADM_CONFIG)
		networkInit(config)
	},
}

func networkInit(config utils.Configuration) {
	file := filepath.Join(utils.CONF_DIR, "flannel.yaml")
	log.Printf("Pod network %s\n", config.MasterConfiguration.Networking.PodSubnet)
	utils.ReplaceString(file, utils.DEFAULT_POD_NETWORK, config.MasterConfiguration.Networking.PodSubnet)
	utils.Run(utils.BASE_DIR, "sysctl", "net.bridge.bridge-nf-call-iptables=1")
	utils.Run(utils.BASE_DIR, "kubectl", "--kubeconfig="+"/etc/kubernetes/admin.conf", "apply", "-f", file)
}

/*
func writeConfFiles() {
	masterConfig := kubeadm.MasterConfiguration{}
	masterConfig.ClusterName = "test"
	masterConfig.Networking.PodSubnet = "10.1.0.0/16"
	masterConfig.Networking.ServiceSubnet = "10.2.0.0/16"
	masterConfig.API.AdvertiseAddress = ""
	masterConfig.API.BindPort = 443
	masterConfig.API.ControlPlaneEndpoint = ""
	masterConfig.APIServerCertSANs = []string{"10.2.1.2"}
	masterConfig.Token = "token"
	masterConfig.CertificatesDir = "/tmp/certs"
	masterConfig.Etcd.Endpoints = []string{"http://127.0.0.1:2379"}
	y, err := yaml.Marshal(masterConfig)
	if err != nil {
		log.Fatalf("Could not serialize master configuration object")
	}
	ioutil.WriteFile("/tmp/config", y, 0644)

}
*/

func kubeadmInit(config string) {
	utils.Run(utils.BASE_DIR, "kubeadm", "init", "--config="+config)
}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().String("cfg", "", "Location of configuration file")
}
