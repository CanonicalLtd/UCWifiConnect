package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
)

func checkSsids(c *netman.Client) bool {
	ssidsFile := os.Getenv("SNAP_COMMON") + "/ssids"
	SSIDs, _, _ := c.Ssids()
	//only write SSIDs when found
	if len(SSIDs) > 0 {
		var out string
		for _, ssid := range SSIDs {
			out += strings.TrimSpace(ssid.Ssid) + ","
		}
		out = out[:len(out)-1]
		err := ioutil.WriteFile(ssidsFile, []byte(out), 0644)
		if err != nil {
			fmt.Println("Error writing ssids to ", ssidsFile)
		} else {
			return true
		}
	}
	return false
}

func main() {
	first := true
	c := netman.DefaultClient()
	// add code

	for {
		if first {
			first = false
			//wait one minute on first run to allow wifi connections
			//TODO: time period to be refined
			time.Sleep(10000 * time.Millisecond)
			//time.Sleep(60000 * time.Millisecond)
		}
		//if not connected to external wifi, set to Management Mode
		if !c.ConnectedWifi(c.GetWifiDevices(c.GetDevices())) {
			fmt.Println("==== No network manager wifi connection")
			// if wlan0 not managed by network manager, set to managed
			wifis, _ := c.WifisManaged(c.GetWifiDevices(c.GetDevices()))
			found := false
			for k, _ := range wifis {
				if k == "wlan0" {
					found = true
				}
			}
			if !found {
				c.SetIfaceManaged("wlan0", c.GetWifiDevices(c.GetDevices()))
			}
			// check for ssids, if any found, take steps
			if checkSsids(c) {
				//TODO: start Management Mode http server. requires new pkg funcs
				fmt.Println("==== Start Management Mode http server...")
				//TODO: start wifi-ap
				fmt.Println("==== Start wifi-ap AP...")
			}
		} else { //Set to operational mode
			//TODO create operational mode server
			fmt.Println("==== Start Operational Mode http server...")
		}
		// wait 5 seconds
		time.Sleep(5000 * time.Millisecond)
	}
}
