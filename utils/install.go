package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	FILE_MODE = 0744
	ETC_DIR   = "/etc/systemd/system"
)

func Install(kubernetesVersion, cniVersion, rootDir string) {
	downloadArtifacts(rootDir, kubernetesVersion, cniVersion)
	writeKubeletServiceFiles(rootDir, kubernetesVersion)
	EnableAndStartService("kubelet.service")
}
func download(fileName string, url string, mode os.FileMode) {
	log.Printf("Downloading %s to location %s", url, fileName)
	_, err := os.Stat(fileName)
	if !os.IsNotExist(err) {
		log.Printf("File already exists %s\n", fileName)
		if err := os.Chmod(fileName, mode); err != nil {
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
	if err := os.Chmod(fileName, mode); err != nil {
		log.Fatalf("Failed to set permissions for file %s, with error %v\n", fileName, err)
	}
}

func writeKubeletServiceFiles(rootDir string, kuberneteVersion string) {
	baseURL := fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", kuberneteVersion)
	//kubelet service
	serviceFile := filepath.Join(ETC_DIR, "kubelet.service")
	download(serviceFile, baseURL+"kubelet.service", FILE_MODE)
	ReplaceString(serviceFile, "/usr/bin", rootDir)

	//kubelet service conf
	err := os.MkdirAll(filepath.Join(ETC_DIR, "kubelet.service.d"), FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir with error %v\n", err)
	}
	confFile := filepath.Join(ETC_DIR, "kubelet.service.d", "10-kubeadm.conf")
	download(confFile, baseURL+"10-kubeadm.conf", FILE_MODE)
	ReplaceString(confFile, "/usr/bin", rootDir)
}

func downloadArtifacts(rootDir, kuberneteVersion, cniVersion string) {
	err := os.MkdirAll(rootDir, FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}
	baseURL := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", kuberneteVersion)
	download(filepath.Join(rootDir, "kubectl"), baseURL+"kubectl", FILE_MODE)
	download(filepath.Join(rootDir, "kubeadm"), baseURL+"kubeadm", FILE_MODE)
	download(filepath.Join(rootDir, "kubelet"), baseURL+"kubelet", FILE_MODE)

	err = os.MkdirAll("/opt/cni/bin", FILE_MODE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", rootDir, err)
	}

	baseURL = fmt.Sprintf("https://github.com/containernetworking/plugins/releases/download/%s/cni-plugins-amd64-%s.tgz", cniVersion, cniVersion)
	tmpFile := fmt.Sprintf("/tmp/cni-plugins-amd64-%s.tgz", cniVersion)
	download(tmpFile, baseURL, FILE_MODE)
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
