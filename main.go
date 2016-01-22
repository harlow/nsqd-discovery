package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func setConfig(nsqAddr string, lookupdPort int, lookupdIPs []net.IP) {
	lookupdTCPAddrs := []string{}
	for _, IP := range lookupdIPs {
		addr := fmt.Sprintf("%s:%d", IP, lookupdPort)
		lookupdTCPAddrs = append(lookupdTCPAddrs, addr)
	}

	body, err := json.Marshal(lookupdTCPAddrs)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("setting lookupdTCPAddrs: %s", body)

	configAddr := "http://" + nsqAddr + "/config/nsqlookupd_tcp_addresses"
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
		dnsAddr  = flag.String("lookupd-dns-address", "", "The DNS address of nsqlookupd")
		nsqdAddr = flag.String("nsqd-http-address", "0.0.0.0:4151", "The HTTP address of nsqd")
	)
	flag.Parse()

	if *dnsAddr == "" {
		fmt.Println("Error: required arg -lookupd-dns-address")
		return
	}

	lookupdIPs, err := net.LookupIP(*dnsAddr)
	if err != nil {
		log.Fatalf("net.LookupIP error: %s", err)
	}

	if len(lookupdIPs) == 0 {
		log.Fatalf("no IPs found for %s", *dnsAddr)
	}

	setConfig(*nsqdAddr, *lookupdPort, lookupdIPs)
	ticker := time.Tick(15 * time.Second)

	for {
		select {
		case <-ticker:
			lookupdIPs, err := net.LookupIP(*dnsAddr)
			if err != nil {
				log.Fatal(err)
			}

			if len(lookupdIPs) == 0 {
				log.Printf("No IP addresses found for %s", *dnsAddr)
				continue
			}

			setConfig(*nsqdAddr, *lookupdPort, lookupdIPs)
		}
	}
}
