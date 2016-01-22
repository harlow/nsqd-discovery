package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// test equality between two arrays
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// get the current lookupd IPs from nsqd config
func getConfgIPs(configAddr string) ([]string, error) {
	configIPs := []string{}

	res, err := http.Get(configAddr)
	if err != nil {
		return configIPs, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return configIPs, err
	}

	json.Unmarshal(body, &configIPs)
	return configIPs, nil
}

// get lookupd IPs from DNS
func getLookupdIPs(dnsAddr *string, port *int) ([]string, error) {
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

// set lookupd addrs on nsqd
func setConfig(configAddr string, lookupdAddrs []string) error {
	body, err := json.Marshal(lookupdAddrs)
	if err != nil {
		return err
	}

	log.Printf("type=info func=setConfig ips=%s", body)

	req, err := http.NewRequest("PUT", configAddr, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		log.Printf("type=error func=setConfig status=%d", res.StatusCode)
	}

	return nil
}

func main() {
	var (
		lookupdPort = flag.Int("lookupd-tcp-port", 4160, "The nsqlookupd tcp port")
		dnsAddr     = flag.String("lookupd-dns-address", "", "The DNS address of nsqlookupd")
		nsqdAddr    = flag.String("nsqd-http-address", "0.0.0.0:4151", "The HTTP address of nsqd")
	)
	flag.Parse()

	configAddr := "http://" + *nsqdAddr + "/config/nsqlookupd_tcp_addresses"

	if *dnsAddr == "" {
		log.Fatal("arg -lookupd-dns-address is required")
	}

	lookupdIPs, err := getLookupdIPs(dnsAddr, lookupdPort)
	if err != nil {
		log.Fatalf("type=error func=getLookupdIPs msg=%s", err)
	}

	if len(lookupdIPs) == 0 {
		log.Printf("type=warn msg=%s addr=%s", "no dns records", *dnsAddr)
	}

	err = setConfig(configAddr, lookupdIPs)
	if err != nil {
		log.Printf("type=error func=setConfig msg=%s", err)
	}

	ticker := time.Tick(15 * time.Second)
	for {
		select {
		case <-ticker:
			lookupdIPs, err = getLookupdIPs(dnsAddr, lookupdPort)
			if err != nil {
				log.Printf("type=error func=getLookupdIPs msg=%s", err)
			}

			if len(lookupdIPs) == 0 {
				log.Printf("type=warn msg=%s addr=%s", "no dns records", *dnsAddr)
				continue
			}

			configIPs, err := getConfgIPs(configAddr)
			if err != nil {
				log.Printf("type=error func=getConfgIPs msg=%s", err)
				continue
			}

			if eq(lookupdIPs, configIPs) {
				log.Printf("type=info func=eq msg=%s", "config up to date")
				continue
			}

			err = setConfig(configAddr, lookupdIPs)
			if err != nil {
				log.Printf("type=error func=setConfig msg=%s", err)
				continue
			}
		}
	}
}
