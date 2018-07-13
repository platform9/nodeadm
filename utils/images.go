package utils

import (
	"fmt"

	"github.com/platform9/nodeadm/constants"
)

var DOCKER_IMAGES = []string{
	constants.KEEPALIVED_IMG,
	fmt.Sprintf("k8s.gcr.io/kube-apiserver-amd64:%s", constants.KUBERNETES_VERSION),
	fmt.Sprintf("k8s.gcr.io/kube-controller-manager-amd64:%s", constants.KUBERNETES_VERSION),
	fmt.Sprintf("k8s.gcr.io/kube-scheduler-amd64:%s", constants.KUBERNETES_VERSION),
	fmt.Sprintf("k8s.gcr.io/kube-proxy-amd64:%s", constants.KUBERNETES_VERSION),
	fmt.Sprintf("k8s.gcr.io/k8s-dns-sidecar-amd64:%s", constants.KUBE_DNS_VERSION),
	fmt.Sprintf("k8s.gcr.io/k8s-dns-kube-dns-amd64:%s", constants.KUBE_DNS_VERSION),
	fmt.Sprintf("k8s.gcr.io/k8s-dns-dnsmasq-nanny-amd64:%s", constants.KUBE_DNS_VERSION),
	"quay.io/coreos/flannel:v0.10.0-amd64",
	"k8s.gcr.io/pause-amd64:3.1",
	"metallb/speaker:master",
	"metallb/controller:master",
}

func GetImages() []string {
	return DOCKER_IMAGES
}
