package utils

import (
	"path/filepath"
)

const (
	KUBERNETES_VERSION  = "v1.10.4"
	CNI_VERSION         = "v0.6.0"
	BASE_DIR            = "/opt/bin"
	CNI_BASE_DIR        = "/opt/cni/bin"
	SYSTEMD_DIR         = "/etc/systemd/system"
	CONFIG_DIR          = "conf"
	FLANNEL_VERSION     = "v0.10.0"
	DEFAULT_POD_NETWORK = "10.244.0.0/16"
	DEFAULT_DNS_IP      = "10.96.0.10"
	DEFAULT_ROUTER_ID   = 42
	K8S_IMG_VERSION     = "1.14.7"
	KEEPALIVED_IMG      = "platform9/keepalived:v2.0.4"
)

var KUBE_DIR = filepath.Join(BASE_DIR, "kubernetes-"+KUBERNETES_VERSION)
var CNI_DIR = filepath.Join(CNI_BASE_DIR, "cni-"+CNI_VERSION)
var CONF_DIR = filepath.Join(BASE_DIR, CONFIG_DIR)
