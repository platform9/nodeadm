package apis

import (
	kubeadmv1alpha1 "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
)

type NodeadmConfiguration struct {
	MasterConfiguration kubeadmv1alpha1.MasterConfiguration `json:"masterConfiguration"`
	VIPConfiguration    VIPConfiguration                    `json:"vipConfiguration"`
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
	NetworkInterface string `json:"networkInterface"`
}
