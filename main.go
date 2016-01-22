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
func getConfgIPs(configAddr string) []string {
	res, err := http.Get(configAddr)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	configIPs := []string{}
	json.Unmarshal(body, &configIPs)

	return configIPs
}

// get lookupd IPs from DNS
func getLookupdIPs(dnsAddr *string, port *int) []string {
	addrs := []string{}

	IPs, err := net.LookupIP(*dnsAddr)
	if err != nil {
		// log warning...
		return addrs
	}

	for _, IP := range IPs {
		addr := fmt.Sprintf("%s:%d", IP, *port)
		addrs = append(addrs, addr)
	}

	return addrs
}

// set lookupd addrs on nsqd
func setConfig(configAddr string, lookupdAddrs []string) {
	body, err := json.Marshal(lookupdAddrs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("setting lookupd addresses to %s", body)

	req, err := http.NewRequest("PUT", configAddr, bytes.NewBuffer(body))
	if err != nil {
		log.Fatalf("http.NewRequest error: %s", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("client.Do error: %s", err)
	}

	if res.StatusCode != 200 {
		log.Printf("nsqd responded with status: %d", res.StatusCode)
	}
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

	lookupdIPs := getLookupdIPs(dnsAddr, lookupdPort)
	if len(lookupdIPs) == 0 {
		log.Fatalf("no IPs found for %s", *dnsAddr)
	}

	setConfig(configAddr, lookupdIPs)
	ticker := time.Tick(15 * time.Second)
	for {
		select {
		case <-ticker:
			lookupdIPs = getLookupdIPs(dnsAddr, lookupdPort)
			if len(lookupdIPs) == 0 {
				log.Printf("no addresses found for %s", *dnsAddr)
				continue
			}

			configIPs := getConfgIPs(configAddr)
			if eq(lookupdIPs, configIPs) {
				log.Println("nsqd is in sync with dns addresses")
				continue
			}

			setConfig(configAddr, lookupdIPs)
		}
	}
}
