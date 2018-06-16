package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
)

// nodeCmd represents the cluster command
var nodeCmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initalize the master node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var rootDir = filepath.Join(utils.BASE_DIR, utils.KUBERNETES_VERSION)
		var confDir = filepath.Join(utils.BASE_DIR, utils.CONFIG_DIR)
		utils.Install(utils.KUBERNETES_VERSION, utils.CNI_VERSION, rootDir)
		cfgFile := cmd.Flag("cfg").Value.String()
		kubeadmInit(cfgFile, rootDir)
		networkInit(confDir, cfgFile, rootDir, utils.FLANNEL_VERSION)
	},
}

func networkInit(confDir, cfgFile, rootDir, flannelVersion string) {
	os.MkdirAll(confDir, utils.FILE_MODE)
	url := fmt.Sprintf("https://raw.githubusercontent.com/coreos/flannel/%s/Documentation/kube-flannel.yml", flannelVersion)
	file := filepath.Join(confDir, "flannel.yaml")
	utils.Download(file, url, utils.FILE_MODE)
	log.Printf("Pod network %s\n", getPodSubnet(cfgFile))
	utils.ReplaceString(file, utils.DEFAULT_POD_NETWORK, getPodSubnet(cfgFile))
	utils.Run(rootDir, "kubectl", "--kubeconfig="+"/etc/kubernetes/admin.conf", "apply", "-f", file)
}

func getPodSubnet(file string) string {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file %s with error %v\n", file, err)
	}
	masterConfig := kubeadm.MasterConfiguration{}
	yaml.Unmarshal(bytes, &masterConfig)
	return masterConfig.Networking.PodSubnet
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

func kubeadmInit(config, rootDir string) {
	utils.Run(rootDir, "kubeadm", "init", "--config="+config)
}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().String("cfg", "", "Location of configuration file")
	nodeCmdInit.Flags().String("vip", "192.168.10.5", "VIP ip to be used for multi master setup")
}
