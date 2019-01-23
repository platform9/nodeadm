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
	fmt.Sprintf("quay.io/coreos/flannel:%s-amd64", constants.FlannelVersion),
	fmt.Sprintf("k8s.gcr.io/pause:%s", constants.PauseContainerVersion),
	fmt.Sprintf("metallb/speaker:%s", constants.MetalLBVersion),
	fmt.Sprintf("metallb/controller:%s", constants.MetalLBVersion),
}

func GetImages() []string {
	return DOCKER_IMAGES
}
