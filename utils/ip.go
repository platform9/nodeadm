package utils

import (
	"log"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

func GetIPFromSubnet(subnet string, hostID int) string {
	_, ipv4Net, err := net.ParseCIDR(subnet)
	if err != nil {
		log.Fatalf("Failed to parse service subnet %s with error %v", subnet, err)
	}
	ip, err := cidr.Host(ipv4Net, hostID)
	if err != nil {
		log.Fatalf("Failed to get ip from subnet %s with error %v", ip, err)
	}
	return ip.String()
}
