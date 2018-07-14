package apis

import (
	"github.com/platform9/nodeadm/constants"
	kubeadmv1alpha1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
)

// SetInitDefaults sets defaults on the configuration used by init
func SetInitDefaults(config *InitConfiguration) {
	kubeadmv1alpha1.SetDefaults_MasterConfiguration(&config.MasterConfiguration)
	config.MasterConfiguration.KubernetesVersion = constants.KUBERNETES_VERSION
	config.MasterConfiguration.NoTaintMaster = true
}

// SetJoinDefaults sets defaults on the configuration used by join
func SetJoinDefaults(config *JoinConfiguration) {
	if config.Networking.ServicesCIDR == "" {
		config.Networking.ServicesCIDR = constants.DefaultServicesCIDR
	}
	if config.Networking.ServiceDomain == "" {
		config.Networking.ServiceDomain = constants.DEFAULT_SERVICE_DOMAIN
	}
}
