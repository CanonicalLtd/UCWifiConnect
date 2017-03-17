package main

import (
	"flag"
	"fmt"

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
	switch {
	case opts.show:
		wifiap.Show()
	case len(opts.ssid) > 1:
		wifiap.SetSsid(opts.ssid)
	case len(opts.passphrase) > 1:
		if len(opts.passphrase) < 13 {
			fmt.Println("Passphrase must be at least 13 chars in length. Please try again.")
			return
		}
		wifiap.SetPassphrase(opts.passphrase)
	case opts.enable:
		wifiap.Enable()
	case opts.disable:
		wifiap.Disable()
	}
}
