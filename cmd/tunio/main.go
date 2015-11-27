package main

import (
	"flag"
	"github.com/getlantern/tunio"
	"log"
)

var (
	deviceName = flag.String("tundev", "tun0", "TUN device name.")
	deviceIP   = flag.String("netif-ipaddr", "", "Address of the virtual router inside the TUN device.")
	deviceMask = flag.String("netif-netmask", "255.255.255.0", "Network mask that defines the traffic that is going to be redirected to the TUN device.")
	proxyAddr  = flag.String("proxy-addr", "", "Lantern proxy address.")
	udpgwAddr  = flag.String("udpgw-remote-server-addr", "", "udpgw remote server address (optional).")
)

func main() {
	flag.Parse()

	if *deviceName == "" {
		log.Fatal("missing tundev")
	}

	if *deviceIP == "" {
		log.Fatal("missing netif-ipaddr")
	}

	if *proxyAddr == "" {
		log.Fatal("missing proxy-addr")
	}

	log.Printf("Configuring device %q (ipaddr: %q, netmask: %q)", *deviceName, *deviceIP, *deviceMask)

	dialer := tunio.NewLanternDialer(*proxyAddr, nil)

	if err := tunio.ConfigureTUN(*deviceName, *deviceIP, *deviceMask, *udpgwAddr, dialer); err != nil {
		log.Fatalf("Failed to configure device: %q", err)
	}
}
