package apis

import (
	"fmt"

	log "github.com/platform9/nodeadm/pkg/logrus"

	"github.com/Jeffail/gabs"
)

// ValidateInit validates the configuration used by the init verb
func ValidateInit(config *InitConfiguration) []error {
	var errorList []error
	if err := vipAndControlPlaneEndpointAreEqual(config); err != nil {
		errorList = append(errorList, err)
	}
	if err := podSubnetDefined(config); err != nil {
		errorList = append(errorList, err)
	}
	if err := kubernetesVersionDefined(config); err != nil {
		errorList = append(errorList, err)
	}
	return errorList
}

// ValidateJoin validates the configuration used by the join verb
func ValidateJoin(config *JoinConfiguration) []error {
	var errorList []error
	if err := tokenDefined(config); err != nil {
		errorList = append(errorList, err)
	}
	if err := discoveryTokenAPIServersDefined(config); err != nil {
		errorList = append(errorList, err)
	}
	if err := discoveryTokenCACertHashesDefined(config); err != nil {
		errorList = append(errorList, err)
	}
	return errorList
}

// vipAndControlPlaneEndpointAreEqual checks that vipConfiguration.IP and
// masterConfiguration.api.controlPlaneEndpoint are equal, if
// vipConfiguration.IP is defined.
func vipAndControlPlaneEndpointAreEqual(config *InitConfiguration) error {
	if config.VIPConfiguration == nil {
		return nil
	}
	p, err := gabs.Consume(config.MasterConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse masterConfiguration: %s", err)
	}

	if !p.ExistsP("api.controlPlaneEndpoint") {
		return fmt.Errorf("masterConfiguration.api.controlPlaneEndpoint must be equal to vipConfiguration.IP")
	}
	cep, ok := p.Path("api.controlPlaneEndpoint").Data().(string)
	if !ok {
		return fmt.Errorf("masterConfiguration.api.controlPlaneEndpoint must be a string")
	}
	if cep != config.VIPConfiguration.IP {
		return fmt.Errorf("masterConfiguration.api.controlPlaneEndpoint must be equal to vipConfiguration.IP")
	}
	return nil
}

// podSubnetDefined checks that masterConfiguration.networking.podSubnet is
// defined
func podSubnetDefined(config *InitConfiguration) error {
	p, err := gabs.Consume(config.MasterConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse masterConfiguration: %s", err)
	}
	if !p.ExistsP("networking.podSubnet") {
		return fmt.Errorf("masterConfiguration.networking.podSubnet must be defined for flannel to work")
	}
	return nil
}

// kubernetesVersionDefined checks that masterConfiguration.kubernetesVersion is
// defined
func kubernetesVersionDefined(config *InitConfiguration) error {
	p, err := gabs.Consume(config.MasterConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse masterConfiguration: %s", err)
	}
	if !p.ExistsP("kubernetesVersion") {
		log.Warn("masterConfiguration.kubernetesVersion must be defined when internet egress is not available")
	}
	return nil
}

// tokenDefined checks that nodeConfiguration.token is defined
func tokenDefined(config *JoinConfiguration) error {
	p, err := gabs.Consume(config.NodeConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse nodeConfiguration: %s", err)
	}
	if !p.ExistsP("token") {
		return fmt.Errorf("nodeConfiguration.token must be defined")
	}
	return nil
}

// discoveryTokenCACertHashesDefined checks that
// nodeConfiguration.discoveryTokenCACertHashes has at least non-empty item
func discoveryTokenCACertHashesDefined(config *JoinConfiguration) error {
	p, err := gabs.Consume(config.NodeConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse nodeConfiguration: %s", err)
	}
	if !p.ExistsP("discoveryTokenCACertHashes") {
		return fmt.Errorf("nodeConfiguration.discoveryTokenCACertHashes must be defined")
	}
	hashkey := p.Path("discoveryTokenCACertHashes")
	children, err := hashkey.Children()
	if err != nil {
		return fmt.Errorf("unable to parse nodeConfiguration.discoveryTokenCACertHashes")
	}
	if len(children) == 0 {
		return fmt.Errorf("nodeConfiguration.discoveryTokenCACertHashes array must have at least one item")
	}
	for i, child := range children {
		hash, ok := child.Data().(string)
		if !ok {
			return fmt.Errorf("nodeConfiguration.discoveryTokenCACertHashes[%d] must be a string", i)
		}
		if hash == "" {
			return fmt.Errorf("nodeConfiguration.discoveryTokenCACertHashes[%d] is an empty string", i)
		}
	}
	return nil
}

// discoveryTokenAPIServersDefined checks that
// nodeConfiguration.discoveryTokenAPIServers has at least non-empty item
func discoveryTokenAPIServersDefined(config *JoinConfiguration) error {
	p, err := gabs.Consume(config.NodeConfiguration)
	if err != nil {
		return fmt.Errorf("unable to parse nodeConfiguration: %s", err)
	}
	if !p.ExistsP("discoveryTokenAPIServers") {
		return fmt.Errorf("nodeConfiguration.discoveryTokenAPIServers must be defined")
	}
	hashkey := p.Path("discoveryTokenAPIServers")
	children, err := hashkey.Children()
	if err != nil {
		return fmt.Errorf("unable to parse nodeConfiguration.discoveryTokenAPIServers")
	}
	if len(children) == 0 {
		return fmt.Errorf("nodeConfiguration.discoveryTokenAPIServers array must have at least one item")
	}
	for i, child := range children {
		hash, ok := child.Data().(string)
		if !ok {
			return fmt.Errorf("nodeConfiguration.discoveryTokenAPIServers[%d] must be a string", i)
		}
		if hash == "" {
			return fmt.Errorf("nodeConfiguration.discoveryTokenAPIServers[%d] is an empty string", i)
		}
	}
	return nil
}
