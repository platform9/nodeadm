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
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
)

// nodeCmd represents the cluster command
var nodeCmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initalize the master node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		rootDir := filepath.Join(utils.BASE_DIR, utils.KUBERNETES_VERSION)
		confDir := filepath.Join(utils.BASE_DIR, utils.CONFIG_DIR)
		file := cmd.Flag("cfg").Value.String()
		bytes, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Failed to read file %s with error %v\n", file, err)
		}
		masterConfig := kubeadm.MasterConfiguration{}
		yaml.Unmarshal(bytes, &masterConfig)
		masterConfig.KubernetesVersion = utils.KUBERNETES_VERSION
		/*
			kubeadm.SetDefaults_MasterConfiguration(&masterConfig)
			masterConfig.KubeletConfiguration.BaseConfig.ClusterDNS = []string{utils.GetIPFromSubnet(masterConfig.Networking.ServiceSubnet, 10)}
			bytes, err = yaml.Marshal(masterConfig)
			if err != nil {
				log.Fatalf("Failed to marshal master config with err %v\n", err)
			}
			err = ioutil.WriteFile(file, bytes, utils.FILE_MODE)
			if err != nil {
				log.Fatalf("Failed to write file %s with error %v\n", file, err)
			}
		*/
		routerId := cmd.Flag("routerId").Value.String()
		intf := cmd.Flag("interface").Value.String()
		vip := cmd.Flag("vip").Value.String()
		utils.InstallMasterComponents(rootDir, routerId, intf, vip, &masterConfig)
		kubeadmInit(file, rootDir)
		networkInit(confDir, cfgFile, rootDir, utils.FLANNEL_VERSION, masterConfig)
	},
}

func networkInit(confDir, cfgFile, rootDir, flannelVersion string, masterConfig kubeadm.MasterConfiguration) {
	os.MkdirAll(confDir, utils.FILE_MODE)
	url := fmt.Sprintf("https://raw.githubusercontent.com/coreos/flannel/%s/Documentation/kube-flannel.yml", flannelVersion)
	file := filepath.Join(confDir, "flannel.yaml")
	utils.Download(file, url, utils.FILE_MODE)
	log.Printf("Pod network %s\n", masterConfig.Networking.PodSubnet)
	utils.ReplaceString(file, utils.DEFAULT_POD_NETWORK, masterConfig.Networking.PodSubnet)
	utils.Run(rootDir, "sysctl", "net.bridge.bridge-nf-call-iptables=1")
	utils.Run(rootDir, "kubectl", "--kubeconfig="+"/etc/kubernetes/admin.conf", "apply", "-f", file)
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
	nodeCmdInit.Flags().String("routerId", "42", "id of router to be used for keepalived")
	nodeCmdInit.Flags().String("interface", "eth0", "interface used for keepalived")
}
