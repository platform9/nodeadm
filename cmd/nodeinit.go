package cmd

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/platform9/nodeadm/utils"
	"github.com/spf13/cobra"
)

const (
	KUBERNETES_VERSION = "v1.9.6"
	CNI_VERSION        = "v0.6.0"
)

// nodeCmd represents the cluster command
var nodeCmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initalize the master node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		var rootDir = filepath.Join("/opt/bin/", KUBERNETES_VERSION)
		utils.Install(KUBERNETES_VERSION, CNI_VERSION, rootDir)
		kubeadmInit(cmd.Flag("cfg").Value.String(), rootDir)
	},
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
	currentPath := os.Getenv("PATH")
	os.Setenv("PATH", currentPath+":"+rootDir)
	log.Printf("Updated PATH variable = %s", os.Getenv("PATH"))
	log.Printf("Running command %s %s %s", filepath.Join(rootDir, "kubeadm"), "init", "--config="+config)

	cmd := exec.Command(filepath.Join(rootDir, "kubeadm"), "init", "--config="+config)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to run command %s with error %v\n", "kubeadm init", err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Failed to get output of command %s with error %v\n", "kubeadm init", err)
	}

}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().String("cfg", "", "Location of configuration file")
	nodeCmdInit.Flags().String("vip", "192.168.10.5", "VIP ip to be used for multi master setup")
}
