package main

import (
	"flag"
	"fmt"

	"github.com/CanonicalLtd/UCWifiConnect/utils"
	"github.com/CanonicalLtd/UCWifiConnect/wifiap"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type options struct {
	show       bool
	enable     bool
	enabled    bool
	disable    bool
	ssid       string
	passphrase string
	security   string
	verbose    bool
	err        string
}

func args() *options {

	opts := &options{}
	flag.BoolVar(&opts.show, "show", false, "Show the wifi-ap configuration")
	flag.BoolVar(&opts.enable, "ap-on", false, "Turn on the AP")
	flag.BoolVar(&opts.enabled, "enabled", false, "Check if the AP is UP")
	flag.BoolVar(&opts.disable, "ap-off", false, "Turn off the AP")
	flag.StringVar(&opts.ssid, "ssid", "", "Set the AP's SSID")
	flag.StringVar(&opts.passphrase, "passphrase", "", "Set the AP's passphrase")
	flag.BoolVar(&opts.verbose, "verbose", false, "Display verbose output")
	opts.security = "wpa2"
	flag.Parse()
	return opts
}

func main() {
	opts := args()
	if len(opts.err) > 0 {
		fmt.Printf("%q. Stopping.\n", opts.err)
		return
	}

	wifiAPClient := wifiap.DefaultClient()
	var err error
	var result map[string]interface{}

	switch {
	case opts.show:
		result, err = wifiAPClient.Show()
		if result != nil {
			utils.PrintMapSorted(result)
			return
		}
	case len(opts.ssid) > 1:
		err = wifiAPClient.SetSsid(opts.ssid)
	case len(opts.passphrase) > 1:
		err = wifiAPClient.SetPassphrase(opts.passphrase)
	case opts.enable:
		err = wifiAPClient.Enable()
	case opts.disable:
		err = wifiAPClient.Disable()
	case opts.enabled:
		res, err := wifiAPClient.Enabled()
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			if res {
				fmt.Println("Wifi-ap is UP")
			} else {
				fmt.Println("Wifi-ap is Down")
			}
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}
