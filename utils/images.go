package utils

import (
	"fmt"

	"github.com/platform9/nodeadm/constants"
)

var DOCKER_IMAGES = []string{
	constants.KeepalivedImg,
	fmt.Sprintf("k8s.gcr.io/kube-apiserver-amd64:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/kube-controller-manager-amd64:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/kube-scheduler-amd64:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/kube-proxy-amd64:%s", constants.KubernetesVersion),
	fmt.Sprintf("k8s.gcr.io/k8s-dns-sidecar-amd64:%s", constants.KubeDNSVersion),
	fmt.Sprintf("k8s.gcr.io/k8s-dns-kube-dns-amd64:%s", constants.KubeDNSVersion),
	fmt.Sprintf("k8s.gcr.io/k8s-dns-dnsmasq-nanny-amd64:%s", constants.KubeDNSVersion),
	"quay.io/coreos/flannel:v0.10.0-amd64",
	"k8s.gcr.io/pause-amd64:3.1",
	"metallb/speaker:master",
	"metallb/controller:master",
}

func GetImages() []string {
	return DOCKER_IMAGES
}
