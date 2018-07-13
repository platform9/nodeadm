package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/platform9/nodeadm/constants"
)

type Artifact struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Upstream string `json:"upstream"`
	Local    string `json:"local"`
}

var NodeArtifact = []Artifact{
	{
		Name:     "kubeadm",
		Type:     "executable",
		Upstream: fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     "kubectl",
		Type:     "executable",
		Upstream: fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     "kubelet",
		Type:     "executable",
		Upstream: fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     "kubelet.service",
		Type:     "regular",
		Upstream: fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     "10-kubeadm.conf",
		Type:     "regular",
		Upstream: fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     fmt.Sprintf("cni-plugins-amd64-%s.tgz", constants.CNI_VERSION),
		Type:     "regular",
		Upstream: fmt.Sprintf("https://github.com/containernetworking/plugins/releases/download/%s/", constants.CNI_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.CNI_DIR_NAME),
	},
	{
		Name:     "kube-flannel.yml",
		Type:     "regular",
		Upstream: fmt.Sprintf("https://raw.githubusercontent.com/coreos/flannel/%s/Documentation/", constants.FLANNEL_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.FLANNEL_DIR_NAME),
	},
}

func PlaceComponentsFromCache() {
	placeKubeComponents()
	placeCNIPlugin()
	placeKubeletServiceFiles()
	placeNetworkConfig()
}

func placeKubeletServiceFiles() {
	//kubelet service
	serviceFile := filepath.Join(constants.SYSTEMD_DIR, "kubelet.service")
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, "kubelet.service"), serviceFile)
	ReplaceString(serviceFile, "/usr/bin", constants.BASE_INSTALL_DIR)

	//kubelet service conf
	err := os.MkdirAll(filepath.Join(constants.SYSTEMD_DIR, "kubelet.service.d"), constants.EXECUTE)
	if err != nil {
		log.Fatalf("Failed to create dir with error %v\n", err)
	}
	confFile := filepath.Join(constants.SYSTEMD_DIR, "kubelet.service.d", "10-kubeadm.conf")
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, "10-kubeadm.conf"), confFile)
	ReplaceString(confFile, "/usr/bin", constants.BASE_INSTALL_DIR)
}

func placeKubeComponents() {
	err := os.MkdirAll(constants.KUBE_VERSION_INSTALL_DIR, constants.EXECUTE)
	if err != nil {
		log.Fatalf("Failed to create dir %s with error %v\n", constants.KUBE_VERSION_INSTALL_DIR, err)
	}
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, "kubectl"), filepath.Join(constants.KUBE_VERSION_INSTALL_DIR, "kubectl"))
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, "kubeadm"), filepath.Join(constants.KUBE_VERSION_INSTALL_DIR, "kubeadm"))
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME, "kubelet"), filepath.Join(constants.KUBE_VERSION_INSTALL_DIR, "kubelet"))
	CreateSymLinks(constants.KUBE_VERSION_INSTALL_DIR, constants.BASE_INSTALL_DIR, true)
}

func placeCNIPlugin() {
	tmpFile := fmt.Sprintf("cni-plugins-amd64-%s.tgz", constants.CNI_VERSION)
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.CNI_DIR_NAME, tmpFile), filepath.Join("/tmp", tmpFile))
	if _, err := os.Stat(constants.CNI_VERSION_INSTALL_DIR); os.IsNotExist(err) {
		err := os.MkdirAll(constants.CNI_VERSION_INSTALL_DIR, constants.EXECUTE)
		if err != nil {
			log.Fatalf("Failed to create dir %s with error %v\n", constants.CNI_VERSION_INSTALL_DIR, err)
		}
		Run("", "tar", "-xvf", filepath.Join("/tmp", tmpFile), "-C", constants.CNI_VERSION_INSTALL_DIR)
		CreateSymLinks(constants.CNI_VERSION_INSTALL_DIR, constants.CNI_BASE_DIR, true)
	}

}

func placeNetworkConfig() {
	os.MkdirAll(constants.CONF_INSTALL_DIR, constants.EXECUTE)
	Run("", "cp", filepath.Join(constants.CACHE_DIR, constants.FLANNEL_DIR_NAME, "kube-flannel.yml"), filepath.Join(constants.CONF_INSTALL_DIR, "flannel.yaml"))
}

func loadAvailableImages(cli *client.Client) {
	os.MkdirAll(constants.IMAGES_CACHE_DIR, constants.EXECUTE)
	files, err := ioutil.ReadDir(constants.IMAGES_CACHE_DIR)
	if err != nil {
		log.Printf("Failed to list files from dir %s skipping loading images with err %v\n", constants.IMAGES_CACHE_DIR, err)
	}
	for _, file := range files {
		Run("", "docker", "load", "-i", filepath.Join(constants.IMAGES_CACHE_DIR, file.Name()))
	}
}

func PopulateCache() {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Failed to create docker client with error %v", err)
	}
	loadAvailableImages(cli)
	for _, image := range GetImages() {
		//first check if image is already in docker cache
		nameFilter := filters.NewArgs()
		nameFilter.Add("reference", image)
		log.Printf("Checking if image %s is available in docker cache", image)
		list, err := cli.ImageList(context.Background(), types.ImageListOptions{
			Filters: nameFilter,
		})
		if err != nil {
			log.Fatalf("Failed to list images with error %v\n", err)
		}
		if len(list) == 0 {
			log.Printf("Trying to pull image %s", image)
			Run("", "docker", "pull", image)
		}
		list, err = cli.ImageList(context.Background(), types.ImageListOptions{
			Filters: nameFilter,
		})
		imageFile := filepath.Join(constants.IMAGES_CACHE_DIR, strings.Replace(list[0].ID, "sha256:", "", -1)+".tar")
		Run("", "docker", "save", image, "-o", imageFile)
	}
	for _, file := range NodeArtifact {
		mode := constants.READ
		if file.Type == "executable" {
			mode = constants.EXECUTE
		}
		os.MkdirAll(file.Local, constants.EXECUTE)
		Download(filepath.Join(file.Local, file.Name), file.Upstream+file.Name, os.FileMode(mode))
	}
}
