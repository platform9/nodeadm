package apis

import (
	"fmt"

	"github.com/platform9/nodeadm/constants"
)

// ValidateInit validates the configuration used by the init verb
func ValidateInit(config *InitConfiguration) []error {
	var errorList []error
	if config.MasterConfiguration.Networking.ServiceSubnet != config.Networking.ServiceSubnet {
		errorList = append(errorList, fmt.Errorf("configuration conflict: Networking.ServiceSubnet=%q, MasterConfiguration.Networking.ServiceSubnet=%q. Values should be identical, or MasterConfiguration.Networking.ServiceSubnet omitted",
			config.Networking.ServiceSubnet, config.MasterConfiguration.Networking.ServiceSubnet))
	}
	if len(config.MasterConfiguration.Networking.PodSubnet) == 0 {
		// Pod subnet was set through MasterConfiguration.ControllerManagerExtraArgs
		value, ok := config.MasterConfiguration.ControllerManagerExtraArgs[constants.ControllerManagerClusterCIDRKey]
		if !ok || value != config.Networking.PodSubnet {
			errorList = append(errorList, fmt.Errorf("configuration conflict: Networking.PodSubnet=%q, MasterConfiguration.ControllerManagerExtraArgs[%q]. Values should be identical, or MasterConfiguration.ControllerManagerExtraArgs[%q] omitted",
				config.Networking.PodSubnet, constants.ControllerManagerClusterCIDRKey, constants.ControllerManagerClusterCIDRKey))
		}
	} else {
		// Pod subnet was set through MasterConfiguration.Networking.PodSubnet
		if config.MasterConfiguration.Networking.PodSubnet != config.Networking.PodSubnet {
			errorList = append(errorList, fmt.Errorf("Configuration conflict: Networking.PodSubnet=%q, MasterConfiguration.Networking.PodSubnet=%q. Values should be identical, or MasterConfiguration.Networking.PodSubnet omitted.",
				config.Networking.PodSubnet, config.MasterConfiguration.Networking.PodSubnet))
		}
	}
	if config.MasterConfiguration.Networking.DNSDomain != config.Networking.DNSDomain {
		errorList = append(errorList, fmt.Errorf("configuration conflict: Networking.DNSDomain=%q, MasterConfiguration.Networking.DNSDomain=%q. Values should be identical, or MasterConfiguration.Networking.DNSDomain omitted",
			config.Networking.DNSDomain, config.MasterConfiguration.Networking.DNSDomain))
	}
	return errorList
}
