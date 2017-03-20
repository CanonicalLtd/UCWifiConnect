package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

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

func printMapSorted(m map[string]interface{}) {
	sortedKeys := make([]string, 0, len(m))
	for key := range m {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)
	for _, k := range sortedKeys {
		fmt.Fprintf(os.Stdout, "%s: %v\n", k, m[k])
	}
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
			printMapSorted(result)
			return
		}
	case len(opts.ssid) > 1:
		err = wifiAPClient.SetSsid(opts.ssid)
	case len(opts.passphrase) > 1:
		if len(opts.passphrase) < 13 {
			fmt.Println("Passphrase must be at least 13 chars in length. Please try again.")
			return
		}
		err = wifiAPClient.SetPassphrase(opts.passphrase)
	case opts.enable:
		err = wifiAPClient.Enable()
	case opts.disable:
		err = wifiAPClient.Disable()
	}

	if err != nil {
		fmt.Println(err)
	}
}
