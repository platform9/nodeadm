package apis

import (
	"github.com/platform9/nodeadm/constants"
)

// SetControllerManagerExtraArgs sets controller manager extra args for a given pod network subnet
func setControllerManagerExtraArgs(config *InitConfiguration) {
	if config.MasterConfiguration.ControllerManagerExtraArgs == nil {
		config.MasterConfiguration.ControllerManagerExtraArgs = make(map[string]string)
	}
	if _, ok := config.MasterConfiguration.ControllerManagerExtraArgs[constants.ControllerManagerAllocateNodeCIDRsKey]; !ok {
		config.MasterConfiguration.ControllerManagerExtraArgs[constants.ControllerManagerAllocateNodeCIDRsKey] = constants.ControllerManagerAllocateNodeCIDRs
	}
	if _, ok := config.MasterConfiguration.ControllerManagerExtraArgs[constants.ControllerManagerClusterCIDRKey]; !ok {
		config.MasterConfiguration.ControllerManagerExtraArgs[constants.ControllerManagerClusterCIDRKey] = config.Networking.PodSubnet
	}
	if _, ok := config.MasterConfiguration.ControllerManagerExtraArgs[constants.ControllerManagerNodeCIDRMaskSizeKey]; !ok {
		config.MasterConfiguration.ControllerManagerExtraArgs[constants.ControllerManagerNodeCIDRMaskSizeKey] = constants.ControllerManagerNodeCIDRMaskSize
	}
}
