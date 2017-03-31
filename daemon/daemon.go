package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
	"github.com/CanonicalLtd/UCWifiConnect/wifiap"
)

func scanSsids(path string, c *netman.Client) bool {
	// set wlan- to managed if needed
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
	SSIDs, _, _ := c.Ssids()
	//only write SSIDs when found
	if len(SSIDs) > 0 {
		var out string
		for _, ssid := range SSIDs {
			out += strings.TrimSpace(ssid.Ssid) + ","
		}
		out = out[:len(out)-1]
		err := ioutil.WriteFile(path, []byte(out), 0644)
		if err != nil {
			fmt.Println("Error writing ssids to ", path)
		} else {
			return true
		}
	}
	return false
}

const (
	unknown = iota
	managed
	operational
)

func main() {
	first := true
	c := netman.DefaultClient()
	cw := wifiap.DefaultClient()
	ssidsPath := os.Getenv("SNAP_COMMON") + "/ssids"
	mode := unknown
	for {
		if first {
			first = false
			//wait time period (TBD) on first run to allow wifi connections
			time.Sleep(10000 * time.Millisecond)
		}
		if c.ConnectedWifi(c.GetWifiDevices(c.GetDevices())) {
			fmt.Println("====  Have network manager wifi connection")
			mode = operational
		}

		//if not connected to external wifi: Management Mode
		if mode != operational {
			fmt.Println("==== No network manager wifi connection")
			if mode == unknown {
				if scanSsids(ssidsPath, c) {
					mode = managed
				} else { // recheck ssids on next iter
					continue
				}
			}
			fmt.Println("==== Management Mode")
			res, err := cw.Enabled()
			if err != nil {
				fmt.Println("====== err:", err)
				continue
			}
			if !res {
				//TODO: start Management Mode http server. requires new pkg funcs
				fmt.Println("==== Start wifi-ap AP...")
				cw.Enable()
			}
			fmt.Println("==== Start Management Mode http server...")
		} else { // Operational mode
			//TODO create operational mode server
			mode = operational
			fmt.Println("==== Start Operational Mode http server...")
		}
		// wait 5 seconds
		time.Sleep(5000 * time.Millisecond)
	}
}
