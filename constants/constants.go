package constants

import (
	"fmt"
	"path/filepath"

	netutil "k8s.io/apimachinery/pkg/util/net"
)

const (
	KubernetesVersion                     = "v1.11.9"
	CNIVersion                            = "v0.6.0"
	BaseInstallDir                        = "/opt/bin"
	CNIBaseDir                            = "/opt/cni/bin"
	CNIConfigDir                          = "/etc/cni"
	CNIStateDir                           = "/var/lib/cni"
	SystemdDir                            = "/etc/systemd/system"
	FlannelVersion                        = "v0.10.0"
	MetalLBVersion                        = "master"
	DefaultPodNetwork                     = "10.244.0.0/16"
	DefaultDNSIP                          = "10.96.0.10"
	DefaultServiceSubnet                  = "10.96.0.0/12"
	DefaultDNSDomain                      = "cluster.local"
	DefaultRouterID                       = 42
	KubeadmConfig                         = "/tmp/kubeadm.yaml"
	CoreDNSVersion                        = "1.1.3"
	PauseContainerVersion                 = "3.1"
	KeepalivedImage                       = "platform9/keepalived:v2.0.4"
	CacheDir                              = "/var/cache/nodeadm/"
	Execute                               = 0744
	Read                                  = 0644
	FeatureGates                          = "ExperimentalCriticalPodAnnotation=true"
	Sysctl                                = "/sbin/sysctl"
	ControllerManagerAllocateNodeCIDRsKey = "allocate-node-cidrs"
	ControllerManagerClusterCIDRKey       = "cluster-cidr"
	ControllerManagerNodeCIDRMaskSizeKey  = "node-cidr-mask-size"
	// TODO(puneet) remove when we move to 1.11.
	// Currently set it similar to upstream
	// https://github.com/kubernetes/kubernetes/blob/v1.10.4/cmd/kubeadm/app/phases/controlplane/manifests.go#L281
	ControllerManagerNodeCIDRMaskSize = "24"
	// TODO(puneet) remove when we move to 1.11.
	// Currently set it similar to upstream
	// https://github.com/kubernetes/kubernetes/blob/v1.10.4/cmd/kubeadm/app/phases/controlplane/manifests.go#L340
	ControllerManagerAllocateNodeCIDRs = "true"
	KubeletConfigKubeReservedCPUKey    = "cpu"
	DefaultAPIBindPort                 = 6443
)

const (
	VRRPScriptInterval = 10
	VRRPScriptRise     = 2
	VRRPScriptFall     = 6
	WgetTimeout        = 8
)

var KubeDirName = filepath.Join("kubernetes", KubernetesVersion)
var FlannelDirName = filepath.Join("flannel", FlannelVersion)
var CNIDirName = filepath.Join("cni", CNIVersion)
var CniVersionInstallDir = filepath.Join(CNIBaseDir, CNIVersion)
var ImagesCacheDir = filepath.Join(CacheDir, "images")

const (
	KubeadmFilename                     = "kubeadm"
	KubectlFilename                     = "kubectl"
	KubeletFilename                     = "kubelet"
	KubeletSystemdUnitFilename          = "kubelet.service"
	KubeadmKubeletSystemdDropinFilename = "10-kubeadm.conf"
	FlannelManifestFilename             = "kube-flannel.yml"
	AdminKubeconfigFile                 = "/etc/kubernetes/admin.conf"
	KeepalivedConfigFilename            = "/etc/keepalived/keepalived.conf"
)

var CNIPluginsFilename = fmt.Sprintf("cni-plugins-amd64-%s.tgz", CNIVersion)

func GetHostnameOverride() (string, error) {
	defaultIP, err := netutil.ChooseHostInterface()
	if err != nil {
		return "", fmt.Errorf("failed to get default interface with err %v", err)
	}
	return defaultIP.String(), nil
}
