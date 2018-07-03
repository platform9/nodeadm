package utils

import (
	"path/filepath"
)

const (
	KUBERNETES_VERSION  = "v1.10.4"
	CNI_VERSION         = "v0.6.0"
	BASE_INSTALL_DIR    = "/opt/bin"
	CNI_BASE_DIR        = "/opt/cni/bin"
	SYSTEMD_DIR         = "/etc/systemd/system"
	CONFIG_DIR          = "conf"
	FLANNEL_VERSION     = "v0.10.0"
	DEFAULT_POD_NETWORK = "10.244.0.0/16"
	DEFAULT_DNS_IP      = "10.96.0.10"
	DEFAULT_ROUTER_ID   = 42
	KUBEADM_CONFIG      = "/tmp/kubeadm.yaml"
	KUBE_DNS_VERSION    = "1.14.8"
	KEEPALIVED_IMG      = "platform9/keepalived:v2.0.4"
	CACHE_BASE_DIR      = "/var/cache/"
)

var KUBE_DIR_NAME = "kubernetes-" + KUBERNETES_VERSION
var CNI_DIR_NAME = "cni-" + CNI_VERSION
var FLANNEL_DIR_NAME = "flannel-" + FLANNEL_VERSION
var NODEADM_DIR_NAME = "noedadm-" + KUBERNETES_VERSION

var KUBE_VERSION_INSTALL_DIR = filepath.Join(BASE_INSTALL_DIR, KUBE_DIR_NAME)
var CNI_VERSION_INSTALL_DIR = filepath.Join(CNI_BASE_DIR, CNI_DIR_NAME)
var CONF_INSTALL_DIR = filepath.Join(BASE_INSTALL_DIR, CONFIG_DIR)
var CACHE_DIR = filepath.Join(CACHE_BASE_DIR, NODEADM_DIR_NAME)
var IMAGES_CACHE_DIR = filepath.Join(CACHE_DIR, "images")
