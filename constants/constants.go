package constants

import (
	"fmt"
	"path/filepath"

	netutil "k8s.io/apimachinery/pkg/util/net"
)

const (
	KubernetesVersion                     = "v1.10.4"
	CNIVersion                            = "v0.6.0"
	BaseInstallDir                        = "/opt/bin"
	CNIBaseDir                            = "/opt/cni/bin"
	CNIConfigDir                          = "/etc/cni"
	CNIStateDir                           = "/var/lib/cni"
	SystemdDir                            = "/etc/systemd/system"
	FlannelVersion                        = "v0.10.0"
	DefaultPodNetwork                     = "10.244.0.0/16"
	DefaultDNSIP                          = "10.96.0.10"
	DefaultServiceSubnet                  = "10.96.0.0/12"
	DefaultDNSDomain                      = "cluster.local"
	DefaultRouterID                       = 42
	KubeadmConfig                         = "/tmp/kubeadm.yaml"
	KubeDNSVersion                        = "1.14.8"
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

const (
	// TODO(dlipovetsky) Move fields to configuration
	KubeletEvictionHard                 = "memory.available<600Mi,nodefs.available<10%"
	NodeadmKubeletSystemdDropinFilename = "20-nodeadm.conf"
	NodeadmKubeletSystemdDropinTemplate = `[Service]
Environment="KUBELET_DNS_ARGS=--cluster-dns={{ .ClusterDNS }} --cluster-domain={{ .ClusterDomain }}"
Environment="KUBELET_EXTRA_ARGS=--resolv-conf=/run/systemd/resolve/resolv.conf --max-pods={{ .MaxPods }} --fail-swap-on={{ .FailSwapOn }} --hostname-override={{ .HostnameOverride }} --kube-api-qps={{ .KubeAPIQPS }} --kube-api-burst={{ .KubeAPIBurst }} --feature-gates={{ .FeatureGates}} --eviction-hard={{ .EvictionHard }} --cpu-manager-policy={{ .CPUManagerPolicy }} --kube-reserved={{ .KubeReservedCPU }}"
`
)

func GetHostnameOverride() (string, error) {
	defaultIP, err := netutil.ChooseHostInterface()
	if err != nil {
		return "", fmt.Errorf("failed to get default interface with err %v", err)
	}
	return defaultIP.String(), nil
}
