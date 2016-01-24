package dnscfg

import (
	"fmt"
	"net"
)

// Get the IP addresses from DNS records
// Note: DNS must have A-Record with IP addresses of nsqlookup instances
func Get(dnsAddr *string, port *int) ([]string, error) {
	addrs := []string{}

	ips, err := net.LookupIP(*dnsAddr)
	if err != nil {
		return addrs, err
	}

	for _, ip := range ips {
		addr := fmt.Sprintf("%s:%d", ip, *port)
		addrs = append(addrs, addr)
	}

	return addrs, nil
}
