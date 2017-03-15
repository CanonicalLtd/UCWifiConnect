package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/CanonicalLtd/UCWifiConnect/netman"

)

type options struct {
	getSsids bool
}

func args() *options {
	opts := &options{}
	flag.BoolVar(&opts.getSsids, "get-ssids", false, "Only display SSIDs (don't connect)")
	flag.Parse()
	return opts
}

func main() {
	opts := args()
	SSIDs, ap2device, ssid2ap := netman.Ssids()
	if opts.getSsids {
		var out string
		for _, ssid := range SSIDs {
			out += strings.TrimSpace(ssid.Ssid) + ","
		}
		if  len(out) > 0 {
		       fmt.Printf("%s\n", out[:len(out)-1])
		}
		return
	}
	for _, ssid := range SSIDs {
		fmt.Printf("    %v\n", ssid.Ssid)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Connect to AP. Enter SSID: ")
	ssid, _ := reader.ReadString('\n')
	ssid = strings.TrimSpace(ssid)
	fmt.Print("Enter phasprase:")
	pw, _ := reader.ReadString('\n')
	pw = strings.TrimSpace(pw)

	netman.ConnectAp(ssid, pw, ap2device, ssid2ap)

	return
}
