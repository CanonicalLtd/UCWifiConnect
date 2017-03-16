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
	getSsids           bool
	checkConnectedWifi bool
	disconnectWifi     bool
}

func args() *options {
	opts := &options{}
	flag.BoolVar(&opts.getSsids, "get-ssids", false, "Only display SSIDs (don't connect)")
	flag.BoolVar(&opts.checkConnectedWifi, "check-connected-wifi", false, "Check if connected to external wifi")
	flag.BoolVar(&opts.disconnectWifi, "disconnect-wifi", false, "Disconnect from any and all external wifi")
	flag.Parse()
	return opts
}

func main() {
	opts := args()
	if opts.getSsids {
		SSIDs, _, _ := netman.Ssids()
		var out string
		for _, ssid := range SSIDs {
			out += strings.TrimSpace(ssid.Ssid) + ","
		}
		if len(out) > 0 {
			fmt.Printf("%s\n", out[:len(out)-1])
		}
		return
	}
	if opts.checkConnectedWifi {
		if netman.ConnectedWifi() {
			fmt.Println("Device is connected to external wifi AP")
		} else {
			fmt.Println("Device is not connected to external wifi AP")

		}
		return
	}
	if opts.disconnectWifi {
		netman.DisconnectWifi()
		return
	} else { //connect
		SSIDs, ap2device, ssid2ap := netman.Ssids()

		for _, ssid := range SSIDs {
			fmt.Printf("    %v\n", ssid.Ssid)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Connect to AP. Enter SSID: ")
		ssid, _ := reader.ReadString('\n')
		ssid = strings.TrimSpace(ssid)
		fmt.Print("Enter phasprase: ")
		pw, _ := reader.ReadString('\n')
		pw = strings.TrimSpace(pw)

		netman.ConnectAp(ssid, pw, ap2device, ssid2ap)
	}
	return
}
