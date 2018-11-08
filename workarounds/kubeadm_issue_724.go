package workarounds

import (
	"github.com/platform9/nodeadm/apis"
	"github.com/platform9/nodeadm/constants"
)

// SetControllerManagerExtraArgs sets controller manager extra args for a given pod network subnet
func SetControllerManagerExtraArgs(config *apis.InitConfiguration) {
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
