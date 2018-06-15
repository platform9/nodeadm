package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/platform9/nodeadm/service"
	"github.com/spf13/cobra"
)

const (
	KUBERNETES_VERSION = "v1.9.6"
	FILE_MODE          = 0744
	ETC_DIR            = "/etc/systemd/system"
	CNI_VERSION        = "v0.6.0"
)

var rootDir = filepath.Join("/opt/bin/", KUBERNETES_VERSION)

// nodeCmd represents the cluster command
var nodeCmdInit = &cobra.Command{
	Use:   "init",
	Short: "Initalize the master node with given configuration",
	Run: func(cmd *cobra.Command, args []string) {
		downloadArtifacts()
		writeKubeletServiceFiles()
		service.EnableAndStartService("kubelet.service")
		kubeadmInit(cmd.Flag("cfg").Value.String())
	},
}

func downloadArtifacts() {
	err := os.MkdirAll(rootDir, FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}
	baseURL := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", KUBERNETES_VERSION)
	download(filepath.Join(rootDir, "kubectl"), baseURL+"kubectl")
	download(filepath.Join(rootDir, "kubeadm"), baseURL+"kubeadm")
	download(filepath.Join(rootDir, "kubelet"), baseURL+"kubelet")

	err = os.MkdirAll("/opt/cni/bin", FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}

	baseURL = fmt.Sprintf("https://github.com/containernetworking/plugins/releases/download/%s/cni-plugins-amd64-%s.tgz", CNI_VERSION, CNI_VERSION)
	tmpFile := fmt.Sprintf("/tmp/cni-plugins-amd64-%s.tgz", CNI_VERSION)
	download(tmpFile, baseURL)
	cmd := exec.Command("tar", "-xvf", tmpFile, "-C", "/opt/cni/bin")
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Failed to untar %s with error %v\n", tmpFile, err)
	}

}

func writeKubeletServiceFiles() {
	baseURL := fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", KUBERNETES_VERSION)
	//kubelet service
	serviceFile := filepath.Join(ETC_DIR, "kubelet.service")
	download(serviceFile, baseURL+"kubelet.service")
	replaceString(serviceFile, "/usr/bin", rootDir)

	//kubelet service conf
	err := os.MkdirAll(filepath.Join(ETC_DIR, "kubelet.service.d"), FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir with error %v\n", err)
	}
	confFile := filepath.Join(ETC_DIR, "kubelet.service.d", "10-kubeadm.conf")
	download(confFile, baseURL+"10-kubeadm.conf")
	replaceString(confFile, "/usr/bin", rootDir)
}

func replaceString(file string, from string, to string) {
	read, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Failed to read file %s with error %v", file, err)
	}
	newContents := strings.Replace(string(read), from, to, -1)
	err = ioutil.WriteFile(file, []byte(newContents), 0)
	if err != nil {
		log.Fatalf("Failed to write file %s with error %v", file, err)
	}
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

func download(fileName string, url string) {
	log.Printf("Downloading %s to location %s", url, fileName)
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		log.Printf("File already exists %s\n", fileName)
		if err := os.Chmod(fileName, FILE_MODE); err != nil {
			log.Fatalf("Failed to set permissions for file %s, with error %v\n", fileName, err)
		}
	} else {
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("Failed to create file %s with err %v\n", fileName, err)
		}
		defer file.Close()
		response, err := http.Get(url)
		if err != nil {
			log.Fatalf("Failed to download %s with error %v\n", url, err)
		}
		defer response.Body.Close()
		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatalf("Failed to save file %s with error %v\n", fileName, err)
		}
	}
	if err := os.Chmod(fileName, FILE_MODE); err != nil {
		log.Fatalf("Failed to set permissions for file %s, with error %v\n", fileName, err)
	}

}

func init() {
	rootCmd.AddCommand(nodeCmdInit)
	nodeCmdInit.Flags().String("cfg", "", "Location of configuration file")
	nodeCmdInit.Flags().String("vip", "192.168.10.5", "VIP ip to be used for multi master setup")
}
