package dnscfg

import(
  "net"
  "fmt"
)

// Get the IP addresses from DNS records
// Note: DNS must have A-Record with IP addresses of nsqlookup instances
func Get(dnsAddr *string, port *int) ([]string, error) {
  addrs := []string{}

  IPs, err := net.LookupIP(*dnsAddr)
  if err != nil {
    return addrs, err
  }

  for _, IP := range IPs {
    addr := fmt.Sprintf("%s:%d", IP, *port)
    addrs = append(addrs, addr)
  }

  return addrs, nil
}
