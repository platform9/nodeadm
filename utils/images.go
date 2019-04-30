package utils

import (
	"fmt"

	"github.com/platform9/nodeadm/constants"
)

var DOCKER_IMAGES = []string{
	constants.KeepalivedImage,
	fmt.Sprintf("k8s.gcr.io/kube-apiserver:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/kube-controller-manager:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/kube-scheduler:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/kube-proxy:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/coredns:%s", constants.CoreDNSVersion),
	fmt.Sprintf("quay.io/coreos/flannel:%s-amd64", constants.FlannelVersion),
	fmt.Sprintf("k8s.gcr.io/pause:%s", constants.PauseContainerVersion),
}

func GetImages() []string {
	return DOCKER_IMAGES
}
