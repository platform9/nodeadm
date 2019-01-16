package utils

import (
	"fmt"

	"github.com/platform9/nodeadm/constants"
)

var DOCKER_IMAGES = []string{
	constants.KeepalivedImage,
	fmt.Sprintf("k8s.gcr.io/kube-apiserver-amd64:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/kube-controller-manager-amd64:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/kube-scheduler-amd64:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/kube-proxy-amd64:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/coredns:%s", constants.CoreDNSVersion),
	"quay.io/coreos/flannel:v0.10.0-amd64",
	fmt.Sprintf("k8s.gcr.io/pause:%s", constants.PauseContainerVersion),
	"metallb/speaker:master",
	"metallb/controller:master",
}

func GetImages() []string {
	return DOCKER_IMAGES
}
