package apis

// InitConfiguration specifies the configuration used by the init command
type InitConfiguration struct {
	MasterConfiguration map[string]interface{} `json:"masterConfiguration,omitempty"`
	VIPConfiguration    *VIPConfiguration      `json:"vipConfiguration,omitempty"`
}

// VIPConfiguration specifies the parameters used to provision a virtual IP
// which API servers advertise and accept requests on.
type VIPConfiguration struct {
	// The virtual IP.
	IP string `json:"ip"`
	// The virtual router ID. Must be in the range [0, 254]. Must be unique within
	// a single L2 network domain.
	RouterID int `json:"routerID"`
	// Network interface chosen to create the virtual IP. If it is not specified,
	// the interface of the default gateway is chosen.
	NetworkInterface string `json:"networkInterface,omitempty"`
}

// JoinConfiguration specifies the configuration used by the join command
type JoinConfiguration struct {
	NodeConfiguration map[string]interface{} `json:"nodeConfiguration,omitempty"`
}
