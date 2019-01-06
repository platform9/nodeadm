package utils

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	log "github.com/platform9/nodeadm/pkg/logrus"

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
		Name:     constants.KubeadmFilename,
		Type:     "executable",
		Upstream: fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", constants.KubernetesVersion),
		Local:    filepath.Join(constants.CacheDir, constants.KubeDirName),
	},
	{
		Name:     constants.KubectlFilename,
		Type:     "executable",
		Upstream: fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", constants.KubernetesVersion),
		Local:    filepath.Join(constants.CacheDir, constants.KubeDirName),
	},
	{
		Name:     constants.KubeletFilename,
		Type:     "executable",
		Upstream: fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/linux/amd64/", constants.KubernetesVersion),
		Local:    filepath.Join(constants.CacheDir, constants.KubeDirName),
	},
	{
		Name:     constants.KubeletSystemdUnitFilename,
		Type:     "regular",
		Upstream: fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", constants.KubernetesVersion),
		Local:    filepath.Join(constants.CacheDir, constants.KubeDirName),
	},
	{
		Name:     constants.KubeadmKubeletSystemdDropinFilename,
		Type:     "regular",
		Upstream: fmt.Sprintf("https://raw.githubusercontent.com/kubernetes/kubernetes/%s/build/debs/", constants.KubernetesVersion),
		Local:    filepath.Join(constants.CacheDir, constants.KubeDirName),
	},
	{
		Name:     constants.CNIPluginsFilename,
		Type:     "regular",
		Upstream: fmt.Sprintf("https://github.com/containernetworking/plugins/releases/download/%s/", constants.CNIVersion),
		Local:    filepath.Join(constants.CacheDir, constants.CNIDirName),
	},
	{
		Name:     constants.FlannelManifestFilename,
		Type:     "regular",
		Upstream: fmt.Sprintf("https://raw.githubusercontent.com/coreos/flannel/%s/Documentation/", constants.FlannelVersion),
		Local:    filepath.Join(constants.CacheDir, constants.FlannelDirName),
	},
}

func loadAvailableImages(cli *client.Client) {
	os.MkdirAll(constants.ImagesCacheDir, constants.Execute)
	files, err := ioutil.ReadDir(constants.ImagesCacheDir)
	if err != nil {
		log.Errorf("\nFailed to list files from dir %s skipping loading images with err %v", constants.ImagesCacheDir, err)
	}
	for _, file := range files {
		cmd := exec.Command("docker", "load", "-i", filepath.Join(constants.ImagesCacheDir, file.Name()))
		err = cmd.Run()
		if err != nil {
			log.Fatalf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
		}
	}
}

func PopulateCache() {
	if os.Getuid() != 0 {
		log.Fatalf("Please run with root privileges.\n")
		return
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatalf("Failed to create docker client with error %v", err)
	}
	loadAvailableImages(cli)
	for _, image := range GetImages() {
		//first check if image is already in docker cache
		nameFilter := filters.NewArgs()
		nameFilter.Add("reference", image)
		log.Infof("Checking if image %s is available in docker cache", image)
		list, err := cli.ImageList(context.Background(), types.ImageListOptions{
			Filters: nameFilter,
		})
		if err != nil {
			log.Fatalf("\nFailed to list images with error %v", err)
		}
		if len(list) == 0 {
			log.Infof("Trying to pull image %s", image)
			cmd := exec.Command("docker", "pull", image)
			err = cmd.Run()
			if err != nil {
				log.Fatalf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
			}

		}
		list, err = cli.ImageList(context.Background(), types.ImageListOptions{
			Filters: nameFilter,
		})
		imageFile := filepath.Join(constants.ImagesCacheDir, strings.Replace(list[0].ID, "sha256:", "", -1)+".tar")
		cmd := exec.Command("docker", "save", image, "-o", imageFile)
		err = cmd.Run()
		if err != nil {
			log.Fatalf("failed to run %q: %s", strings.Join(cmd.Args, " "), err)
		}
	}
	for _, file := range NodeArtifact {
		mode := constants.Read
		if file.Type == "executable" {
			mode = constants.Execute
		}
		os.MkdirAll(file.Local, constants.Execute)
		Download(filepath.Join(file.Local, file.Name), file.Upstream+file.Name, os.FileMode(mode))
	}
}
