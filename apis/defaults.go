package apis

import (
	"fmt"

	"github.com/platform9/nodeadm/constants"
	log "github.com/platform9/nodeadm/pkg/logrus"
	netutil "github.com/platform9/nodeadm/pkg/util/net"

	"github.com/Jeffail/gabs"
)

// SetInitDefaults sets defaults on the configuration used by init
func SetInitDefaults(config *InitConfiguration) {
}

// SetInitDynamicDefaults sets defaults derived at runtime
func SetInitDynamicDefaults(config *InitConfiguration) error {
	if err := setAPIBindPort(config); err != nil {
		return fmt.Errorf("unable to set a default value for masterConfiguration.api.bindPort: %s", err)
	}
	if err := setAPIControlPlaneEndpoint(config); err != nil {
		return fmt.Errorf("unable to set a default value for masterConfiguration.api.controlPlaneEndpoint: %s", err)
	}
	if err := setKeepalivedInterface(config); err != nil {
		return fmt.Errorf("unable to set a default value for vipConfiguraton.networkInterface: %s", err)
	}
	return nil
}

// setAPIBindPort sets the API BindPort, if it us not defined, to a default
// value, because the port must be known to define the keepalived health check
// script. If the VIP is not configured, setAPIBindPort does nothing.
func setAPIBindPort(config *InitConfiguration) error {
	if config.VIPConfiguration == nil {
		return nil
	}
	p, err := gabs.Consume(config.MasterConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse kubeadm MasterConfiguration: %s", err)
	}
	if p.ExistsP("api.bindPort") {
		return nil
	}
	log.Infof("Setting masterConfiguration.api.bindPort to %d", constants.DefaultAPIBindPort)
	_, err = p.SetP(constants.DefaultAPIBindPort, "api.bindPort")
	if err != nil {
		return fmt.Errorf("unable to set kubeadm MasterConfiguration.api.bindPort: %s", err)
	}
	return nil
}

// setAPIControlPlaneEndpoint sets the API ControlPlaneEndpoint to the VIP IP,
// if the VIP is configured. If the VIP is not configured,
// setAPIControlPlaneEndpoint does nothing.
func setAPIControlPlaneEndpoint(config *InitConfiguration) error {
	if config.VIPConfiguration == nil {
		return nil
	}
	p, err := gabs.Consume(config.MasterConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse kubeadm MasterConfiguration: %s", err)
	}
	if p.ExistsP("api.controlPlaneEndpoint") {
		return nil
	}
	log.Infof("Setting masterConfiguration.api.controlPlaneEndpoint equal to vipConfiguration.IP (%q)", config.VIPConfiguration.IP)
	_, err = p.SetP(config.VIPConfiguration.IP, "api.controlPlaneEndpoint")
	if err != nil {
		return fmt.Errorf("unable to set kubeadm MasterConfiguration.api.controlPlaneEndpoint: %s", err)
	}
	return nil
}

// setKeepalivedInterface sets the interface, if it is not defined, to be used
// by keepalived. If the VIP is not configured, setKeepalivedInterface does
// nothing.
func setKeepalivedInterface(config *InitConfiguration) error {
	if config.VIPConfiguration == nil {
		return nil
	}
	if config.VIPConfiguration.NetworkInterface == "" {
		iface, err := netutil.ChooseHostInterface()
		if err != nil {
			return err
		}
		log.Infof("Setting vipConfiguration.networkInterface to %s", iface)
		config.VIPConfiguration.NetworkInterface = iface
	}
	return nil
}
