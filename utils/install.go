package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
)

const (
	FILE_MODE = 0744
	ETC_DIR   = "/etc/systemd/system"
)

func Install(kubernetesVersion, cniVersion, rootDir string, masterConfig *kubeadm.MasterConfiguration) {
	downloadArtifacts(rootDir, kubernetesVersion, cniVersion)
	writeKubeletServiceFiles(rootDir, kubernetesVersion)
	EnableAndStartService("kubelet.service")
	if masterConfig != nil {
		ReplaceString(getKubeletServiceConf(), DEFAULT_DNS_IP, GetIPFromSubnet(masterConfig.Networking.ServiceSubnet, 10))
	}
}

func writeKubeletServiceFiles(rootDir string, kuberneteVersion string) {
	baseURL := fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", kuberneteVersion)
	//kubelet service
	serviceFile := filepath.Join(ETC_DIR, "kubelet.service")
	Download(serviceFile, baseURL+"kubelet.service", FILE_MODE)
	ReplaceString(serviceFile, "/usr/bin", rootDir)

	//kubelet service conf
	err := os.MkdirAll(filepath.Join(ETC_DIR, "kubelet.service.d"), FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir with error %v\n", err)
	}
	confFile := filepath.Join(ETC_DIR, "kubelet.service.d", "10-kubeadm.conf")
	Download(confFile, baseURL+"10-kubeadm.conf", FILE_MODE)
	ReplaceString(confFile, "/usr/bin", rootDir)
}

func getKubeletServiceConf() string {
	return filepath.Join(ETC_DIR, "kubelet.service.d", "10-kubeadm.conf")
}

func downloadArtifacts(rootDir, kuberneteVersion, cniVersion string) {
	err := os.MkdirAll(rootDir, FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}
	baseURL := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", kuberneteVersion)
	Download(filepath.Join(rootDir, "kubectl"), baseURL+"kubectl", FILE_MODE)
	Download(filepath.Join(rootDir, "kubeadm"), baseURL+"kubeadm", FILE_MODE)
	Download(filepath.Join(rootDir, "kubelet"), baseURL+"kubelet", FILE_MODE)

	err = os.MkdirAll("/opt/cni/bin", FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}

	baseURL = fmt.Sprintf("https://github.com/containernetworking/plugins/releases/download/%s/cni-plugins-amd64-%s.tgz", cniVersion, cniVersion)
	tmpFile := fmt.Sprintf("/tmp/cni-plugins-amd64-%s.tgz", cniVersion)
	Download(tmpFile, baseURL, FILE_MODE)
	cmd := exec.Command("tar", "-xvf", tmpFile, "-C", "/opt/cni/bin")
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to untar %s with error %v\n", tmpFile, err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Failed to untar %s with error %v\n", tmpFile, err)
	}
}
