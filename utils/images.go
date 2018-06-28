package utils

import "fmt"

var DOCKER_IMAGES = []string{
	KEEPALIVED_IMG,
	fmt.Sprintf("k8s.gcr.io/kube-apiserver-amd64:%s", KUBERNETES_VERSION),
	fmt.Sprintf("k8s.gcr.io/kube-controller-manager-amd64:%s", KUBERNETES_VERSION),
	fmt.Sprintf("k8s.gcr.io/kube-scheduler-amd64:%s", KUBERNETES_VERSION),
	fmt.Sprintf("k8s.gcr.io/kube-proxy-amd64:%s", KUBERNETES_VERSION),
	fmt.Sprintf("k8s.gcr.io/k8s-dns-sidecar-amd64:%s", K8S_IMG_VERSION),
	fmt.Sprintf("k8s.gcr.io/k8s-dns-kube-dns-amd64:%s", K8S_IMG_VERSION),
	fmt.Sprintf("k8s.gcr.io/k8s-dns-dnsmasq-nanny-amd64:%s", K8S_IMG_VERSION),
	"k8s.gcr.io/pause:3.0",
	"metallb/speaker:master",
	"metallb/controller:master",
}

func GetImages() []string {
	return DOCKER_IMAGES
}
