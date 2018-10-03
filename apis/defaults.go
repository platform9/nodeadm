package apis

import (
	"fmt"

	"github.com/platform9/nodeadm/constants"
	kubeadmv1alpha1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
)

// SetInitDefaults sets defaults on the configuration used by init
func SetInitDefaults(config *InitConfiguration) {
	// First set Networking defaults
	SetNetworkingDefaults(&config.Networking)
	// Second set MasterConfiguration.Networking defaults
	SetMasterConfigurationNetworkingDefaultsWithNetworking(config)
	// Third use the remainder of MasterConfiguration defaults
	kubeadmv1alpha1.SetDefaults_MasterConfiguration(&config.MasterConfiguration)
	config.MasterConfiguration.Kind = "MasterConfiguration"
	config.MasterConfiguration.APIVersion = "kubeadm.k8s.io/v1alpha1"
	config.MasterConfiguration.KubernetesVersion = constants.KubernetesVersion
	config.MasterConfiguration.NoTaintMaster = true

	if config.MasterConfiguration.APIServerExtraArgs == nil {
		config.MasterConfiguration.APIServerExtraArgs = map[string]string{
			"feature-gates": constants.FeatureGates,
		}
	} else {
		prevFeatureGates := config.MasterConfiguration.APIServerExtraArgs["feature-gates"]
		if prevFeatureGates == "" {
			config.MasterConfiguration.APIServerExtraArgs["feature-gates"] = constants.FeatureGates
		} else {
			featureGates := prevFeatureGates + "," + constants.FeatureGates
			config.MasterConfiguration.APIServerExtraArgs["feature-gates"] = featureGates
		}
	}

	if config.MasterConfiguration.ControllerManagerExtraArgs == nil {
		config.MasterConfiguration.ControllerManagerExtraArgs = map[string]string{
			"feature-gates": constants.FeatureGates,
		}
	} else {
		prevFeatureGates := config.MasterConfiguration.ControllerManagerExtraArgs["feature-gates"]
		if prevFeatureGates == "" {

		} else {
			featureGates := prevFeatureGates + "," + constants.FeatureGates
			config.MasterConfiguration.ControllerManagerExtraArgs["feature-gates"] = featureGates
		}
	}

	if config.MasterConfiguration.SchedulerExtraArgs == nil {
		config.MasterConfiguration.SchedulerExtraArgs = map[string]string{
			"feature-gates": constants.FeatureGates,
		}
	} else {
		prevFeatureGates := config.MasterConfiguration.SchedulerExtraArgs["feature-gates"]
		if prevFeatureGates == "" {
			config.MasterConfiguration.SchedulerExtraArgs["feature-gates"] = constants.FeatureGates
		} else {
			featureGates := prevFeatureGates + "," + constants.FeatureGates
			config.MasterConfiguration.SchedulerExtraArgs["feature-gates"] = featureGates
		}
	}

}

// SetInitDynamicDefaults sets defaults derived  at runtime
func SetInitDynamicDefaults(config *InitConfiguration) error {
	nodeName, err := constants.GetHostnameOverride()
	if err != nil {
		return fmt.Errorf("unable to dervice hostname override: %v", err)
	}
	config.MasterConfiguration.NodeName = nodeName
	return nil
}

// SetJoinDefaults sets defaults on the configuration used by join
func SetJoinDefaults(config *JoinConfiguration) {
	SetNetworkingDefaults(&config.Networking)
}

// SetNetworkingDefaults sets defaults for the network configuration
func SetNetworkingDefaults(netConfig *Networking) {
	if netConfig.ServiceSubnet == "" {
		netConfig.ServiceSubnet = constants.DefaultServiceSubnet
	}
	if netConfig.DNSDomain == "" {
		netConfig.DNSDomain = constants.DefaultDNSDomain
	}
}

// SetMasterConfigurationNetworkingDefaultsWithNetworking sets defaults with
// values from the top-level network configuration
func SetMasterConfigurationNetworkingDefaultsWithNetworking(config *InitConfiguration) {
	if config.MasterConfiguration.Networking.ServiceSubnet == "" {
		config.MasterConfiguration.Networking.ServiceSubnet = config.Networking.ServiceSubnet
	}
	if config.MasterConfiguration.Networking.PodSubnet == "" {
		config.MasterConfiguration.Networking.PodSubnet = config.Networking.PodSubnet
	}
	if config.MasterConfiguration.Networking.DNSDomain == "" {
		config.MasterConfiguration.Networking.DNSDomain = config.Networking.DNSDomain
	}
}
