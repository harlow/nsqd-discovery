package main

import (
	"flag"
	"log"
	"time"

	"github.com/harlow/nsqd-discovery/dnscfg"
	"github.com/harlow/nsqd-discovery/httpcfg"
)

var (
	ldPort  = flag.Int("lookupd-tcp-port", 4160, "The nsqlookupd TCP port")
	dnsAddr = flag.String("lookupd-dns-address", "", "The nsqlookupd DNS entry")
	cfgAddr = flag.String("config-http-address", "", "The config address")
)

func main() {
	flag.Parse()

	if *cfgAddr == "" {
		log.Fatal("arg -lookupd-cfg-address is required")
	}

	if *dnsAddr == "" {
		log.Fatal("arg -lookupd-dns-address is required")
	}

	cfgURL := "http://" + *cfgAddr + "/config/nsqlookupd_tcp_addresses"

	IPs, err := dnscfg.Get(dnsAddr, ldPort)
	if err != nil {
		log.Fatalf("type=error msg=%s err=%s", "dns lookup", err)
	}

	if len(IPs) == 0 {
		log.Printf("type=warn msg=%s addr=%s", "no dns records", *dnsAddr)
	}

	err = httpcfg.Set(cfgURL, IPs)
	if err != nil {
		log.Printf("type=error msg=%s err=%s", "set config", err)
	}

	configLoop(cfgURL)
}

// continue looking at dns entry for changes in config
func configLoop(cfgURL string) {
	ticker := time.Tick(15 * time.Second)

	for {
		select {
		case <-ticker:
			newIPs, err := dnscfg.Get(dnsAddr, ldPort)
			if err != nil {
				log.Printf("type=error msg=%s err=%s", "dns lookup", err)
				continue
			}

			if len(newIPs) == 0 {
				log.Printf("type=warn msg=%s addr=%s", "no dns records", *dnsAddr)
				continue
			}

			oldIPs, err := httpcfg.Get(cfgURL)
			if err != nil {
				log.Printf("type=error msg=%s err=%s", "get config", err)
				continue
			}

			if eq(newIPs, oldIPs) {
				log.Printf("type=info msg=%s", "config up to date")
				continue
			}

			err = httpcfg.Set(cfgURL, newIPs)
			if err != nil {
				log.Printf("type=error msg=%s err=%s", "set config", err)
				continue
			}

			log.Printf("type=info msg=%s ips=%s", "set config", newIPs)
		}
	}
}

// test equality of two slices
func eq(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		if !contains(b, v) {
			return false
		}
	}
	return true
}

// test if slice contains string
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
