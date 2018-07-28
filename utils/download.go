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
	"github.com/platform9/nodeadm/deprecated"
)

type Artifact struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Upstream string `json:"upstream"`
	Local    string `json:"local"`
}

var NodeArtifact = []Artifact{
	{
		Name:     constants.KubeadmFilename,
		Type:     "executable",
		Upstream: fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     constants.KubectlFilename,
		Type:     "executable",
		Upstream: fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     constants.KubeletFilename,
		Type:     "executable",
		Upstream: fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     constants.KubeletSystemdUnitFilename,
		Type:     "regular",
		Upstream: fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     constants.KubeadmKubeletSystemdDropinFilename,
		Type:     "regular",
		Upstream: fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", constants.KUBERNETES_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.KUBE_DIR_NAME),
	},
	{
		Name:     constants.CNIPluginsFilename,
		Type:     "regular",
		Upstream: fmt.Sprintf("https://github.com/containernetworking/plugins/releases/download/%s/", constants.CNI_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.CNI_DIR_NAME),
	},
	{
		Name:     constants.FlannelManifestFilename,
		Type:     "regular",
		Upstream: fmt.Sprintf("https://raw.githubusercontent.com/coreos/flannel/%s/Documentation/", constants.FLANNEL_VERSION),
		Local:    filepath.Join(constants.CACHE_DIR, constants.FLANNEL_DIR_NAME),
	},
}

func loadAvailableImages(cli *client.Client) {
	os.MkdirAll(constants.IMAGES_CACHE_DIR, constants.EXECUTE)
	files, err := ioutil.ReadDir(constants.IMAGES_CACHE_DIR)
	if err != nil {
		log.Printf("Failed to list files from dir %s skipping loading images with err %v\n", constants.IMAGES_CACHE_DIR, err)
	}
	for _, file := range files {
		deprecated.Run("", "docker", "load", "-i", filepath.Join(constants.IMAGES_CACHE_DIR, file.Name()))
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
			deprecated.Run("", "docker", "pull", image)
		}
		list, err = cli.ImageList(context.Background(), types.ImageListOptions{
			Filters: nameFilter,
		})
		imageFile := filepath.Join(constants.IMAGES_CACHE_DIR, strings.Replace(list[0].ID, "sha256:", "", -1)+".tar")
		deprecated.Run("", "docker", "save", image, "-o", imageFile)
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
