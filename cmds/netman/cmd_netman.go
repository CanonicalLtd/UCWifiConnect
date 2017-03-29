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
	//devel              bool
	getSsids           bool
	checkConnected     bool
	checkConnectedWifi bool
	disconnectWifi     bool
	wifisManaged       bool
	setIfaceManaged    string
}

func args() *options {
	opts := &options{}
	//flag.BoolVar(&opts.devel, "devel", false, "Test a hard coded devel path")
	flag.BoolVar(&opts.getSsids, "get-ssids", false, "Only display SSIDs (don't connect)")
	flag.BoolVar(&opts.checkConnected, "check-connected", false, "Check if connected at all")
	flag.BoolVar(&opts.checkConnectedWifi, "check-connected-wifi", false, "Check if connected to external wifi")
	flag.BoolVar(&opts.disconnectWifi, "disconnect-wifi", false, "Disconnect from any and all external wifi")
	flag.BoolVar(&opts.wifisManaged, "wifis-managed", false, "Show list of wifi interfaces that are managed by network-manager")
	flag.StringVar(&opts.setIfaceManaged, "manage-iface", "", "Set the specified interface to be managed by network-manager.")
	flag.Parse()
	return opts
}

func main() {
	c := netman.DefaultClient()
	opts := args()
	/*
		if opts.devel {
			dvs := c.GetDevices()
			dwvs := c.GetWifiDevices(dvs)
			for _, d := range dwvs {
				fmt.Println(d)
			}
			return
		}
	*/
	if opts.getSsids {
		SSIDs, _, _ := c.Ssids()
		var out string
		for _, ssid := range SSIDs {
			out += strings.TrimSpace(ssid.Ssid) + ","
		}
		if len(out) > 0 {
			fmt.Printf("%s\n", out[:len(out)-1])
		}
		return
	}
	if opts.checkConnected {
		if c.Connected(c.GetDevices()) {
			fmt.Println("Device is connected ")
		} else {
			fmt.Println("Device is not connected")
		}
		return
	}
	if opts.checkConnectedWifi {
		if c.ConnectedWifi(c.GetWifiDevices(c.GetDevices())) {
			fmt.Println("Device is connected to external wifi AP")
		} else {
			fmt.Println("Device is not connected to external wifi AP")
		}
		return
	}
	if opts.disconnectWifi {
		c.DisconnectWifi(c.GetWifiDevices(c.GetDevices()))
		return
	}
	if len(opts.setIfaceManaged) > 0 {
		c.SetIfaceManaged(opts.setIfaceManaged, c.GetWifiDevices(c.GetDevices()))
		return
	}
	if opts.wifisManaged {
		wifis, err := c.WifisManaged(c.GetWifiDevices(c.GetDevices()))
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range wifis {
			fmt.Printf("%s : %s\n", k, v)
		}
		return
	}
	SSIDs, ap2device, ssid2ap := c.Ssids()

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

	c.ConnectAp(ssid, pw, ap2device, ssid2ap)
}
