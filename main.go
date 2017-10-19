package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/harlow/nsqd-discovery/dnscfg"
	"github.com/harlow/nsqd-discovery/httpcfg"
)

var (
	ldPort      = flag.Int("lookupd-tcp-port", 4160, "The nsqlookupd TCP port")
	dnsAddr     = flag.String("lookupd-dns-address", "", "The nsqlookupd DNS entry")
	cfgAddr     = flag.String("config-http-address", "", "The config address")
	httpAddrCfg = flag.Bool("config-addresses-as-http", false, "Config nsqlookupd http addresses")
)

func main() {
	flag.Parse()

	if *cfgAddr == "" {
		fmt.Println("required flag not provided: --config-http-address")
		os.Exit(1)
	}

	if *dnsAddr == "" {
		fmt.Println("required flag not provided: --lookupd-dns-address")
		os.Exit(1)
	}

	log.SetHandler(text.New(os.Stderr))
	ctx := log.WithFields(log.Fields{
		"cfgAddr": *cfgAddr,
		"dnsAddr": *dnsAddr,
		"ldPort":  *ldPort,
	})

	ips, err := dnscfg.Get(dnsAddr, ldPort)
	if err != nil {
		ctx.WithError(err).Error("dns lookup")
	}

	if len(ips) == 0 {
		ctx.Error("no ip addresses found")
	}

	cfgURL := "http://" + *cfgAddr + "/config/nsqlookupd_tcp_addresses"
	if *httpAddrCfg {
		cfgURL = "http://" + *cfgAddr + "/config/nsqlookupd_http_addresses"
	}

	err = httpcfg.Set(cfgURL, ips)
	if err != nil {
		ctx.WithError(err).Error("setting config")
	} else {
		ctx.WithField("ips", ips).Info("setting config")
	}

	go configLoop(ctx, cfgURL)

	http.ListenAndServe(":6060", nil)
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
				ctx.Error("no ip addresses found")
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
