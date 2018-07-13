package apis

import (
	"github.com/platform9/nodeadm/constants"
	kubeadmv1alpha1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
)

func SetInitDefaults(config *NodeadmConfiguration) {
	kubeadmv1alpha1.SetDefaults_MasterConfiguration(&config.MasterConfiguration)
	config.MasterConfiguration.KubernetesVersion = constants.KUBERNETES_VERSION
	config.MasterConfiguration.NoTaintMaster = true
}

