package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func DownloadKubeComponents(rootDir, version string) {
	err := os.MkdirAll(rootDir, FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}

	//Download kubectl, kubeadm, kubelet if needed
	baseURL := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", version)
	Download(filepath.Join(rootDir, "kubectl"), baseURL+"kubectl", FILE_MODE)
	Download(filepath.Join(rootDir, "kubeadm"), baseURL+"kubeadm", FILE_MODE)
	Download(filepath.Join(rootDir, "kubelet"), baseURL+"kubelet", FILE_MODE)
	CreateSymLinks(KUBE_DIR, BASE_DIR, true)

}

func DownloadCNIPlugin(rootDir, version string) {
	err := os.MkdirAll(rootDir, FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}

	baseURL := fmt.Sprintf("https://github.com/containernetworking/plugins/releases/download/%s/cni-plugins-amd64-%s.tgz", version, version)
	tmpFile := fmt.Sprintf("/tmp/cni-plugins-amd64-%s.tgz", version)
	if _, err := os.Stat(CNI_DIR); os.IsNotExist(err) {
		Download(tmpFile, baseURL, FILE_MODE)
		Run(rootDir, "tar", "-xvf", tmpFile, "-C", rootDir)
		CreateSymLinks(CNI_DIR, CNI_BASE_DIR, true)
	}

}

func DownloadNetworkConfig() {
	os.MkdirAll(CONF_DIR, FILE_MODE)
	url := fmt.Sprintf("https://raw.githubusercontent.com/coreos/flannel/%s/Documentation/kube-flannel.yml", FLANNEL_VERSION)
	file := filepath.Join(CONF_DIR, "flannel.yaml")
	Download(file, url, FILE_MODE)
}

func DownloadArtifacts() {
	DownloadKubeComponents(KUBE_DIR, KUBERNETES_VERSION)
	DownloadCNIPlugin(CNI_DIR, CNI_VERSION)
	DownloadKubeletServiceFiles(KUBERNETES_VERSION)
	DownloadNetworkConfig()
}

func DownloadDockerImages() {
	images := GetImages()
	for _, image := range images {
		Run("", "docker", "pull", image)
	}
}
