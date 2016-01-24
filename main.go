package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/harlow/nsqd-discovery/dnscfg"
	"github.com/harlow/nsqd-discovery/httpcfg"
)

var (
	ldPort  = flag.Int("lookupd-tcp-port", 4160, "The nsqlookupd TCP port")
	dnsAddr = flag.String("lookupd-dns-address", "", "The nsqlookupd DNS entry")
	cfgAddr = flag.String("config-http-address", "", "The config address")
)

func main() {
	log.SetHandler(text.New(os.Stderr))
	flag.Parse()
	ensureRequiredFlags()

	ctx := log.WithFields(log.Fields{
		"cfgAddr": *cfgAddr,
		"dnsAddr": *dnsAddr,
		"ldPort":  *ldPort,
	})

	IPs, err := dnscfg.Get(dnsAddr, ldPort)
	if err != nil {
		ctx.WithError(err).Error("dns lookup")
		os.Exit(1)
	}

	if len(IPs) == 0 {
		ctx.Error("no ips found")
	}

	cfgURL := "http://" + *cfgAddr + "/config/nsqlookupd_tcp_addresses"

	err = httpcfg.Set(cfgURL, IPs)
	if err != nil {
		ctx.WithError(err).Error("setting config")
	} else {
		ctx.WithField("ips", IPs).Info("setting config")
	}

	configLoop(ctx, cfgURL)
}

// make sure the appropriate flags have been passed
func ensureRequiredFlags() {
	if *cfgAddr == "" {
		fmt.Println("required flag not provided: --lookupd-cfg-address")
		os.Exit(1)
	}

	if *dnsAddr == "" {
		fmt.Println("required flag not provided: --lookupd-dns-address")
		os.Exit(1)
	}
}

// continue looking at dns entry for changes in config
func configLoop(ctx log.Interface, cfgURL string) {
	ticker := time.Tick(15 * time.Second)

	for {
		select {
		case <-ticker:
			newIPs, err := dnscfg.Get(dnsAddr, ldPort)
			if err != nil {
				ctx.WithError(err).Error("dns lookup")
				continue
			}

			if len(newIPs) == 0 {
				ctx.Warn("no ips found")
				continue
			}

			oldIPs, err := httpcfg.Get(cfgURL)
			if err != nil {
				ctx.WithError(err).Error("getting config")
				continue
			}

			if eq(newIPs, oldIPs) {
				ctx.Info("config up to date")
				continue
			}

			err = httpcfg.Set(cfgURL, newIPs)
			if err != nil {
				ctx.WithError(err).Error("setting config")
				continue
			}

			ctx.WithField("ips", newIPs).Info("setting config")
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
