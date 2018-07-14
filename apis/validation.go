package apis

import (
	"fmt"
)

// ValidateInit validates the configuration used by the init verb
func ValidateInit(config *InitConfiguration) []error {
	var errorList []error
	if config.MasterConfiguration.Networking.ServiceSubnet != config.Networking.ServiceSubnet {
		errorList = append(errorList, fmt.Errorf("Configuration conflict: Networking.ServiceSubnet=%q, MasterConfiguration.Networking.ServiceSubnet=%q. Values should be identical, or MasterConfiguration.Networking.ServiceSubnet omitted.",
			config.Networking.ServiceSubnet, config.MasterConfiguration.Networking.ServiceSubnet))
	}
	if config.MasterConfiguration.Networking.PodSubnet != config.Networking.PodSubnet {
		errorList = append(errorList, fmt.Errorf("Configuration conflict: Networking.PodSubnet=%q, MasterConfiguration.Networking.PodSubnet=%q. Values should be identical, or MasterConfiguration.Networking.PodSubnet omitted.",
			config.Networking.PodSubnet, config.MasterConfiguration.Networking.PodSubnet))
	}
	if config.MasterConfiguration.Networking.DNSDomain != config.Networking.DNSDomain {
		errorList = append(errorList, fmt.Errorf("Configuration conflict: Networking.DNSDomain=%q, MasterConfiguration.Networking.DNSDomain=%q. Values should be identical, or MasterConfiguration.Networking.DNSDomain omitted.",
			config.Networking.DNSDomain, config.MasterConfiguration.Networking.DNSDomain))
	}
	return errorList
}
