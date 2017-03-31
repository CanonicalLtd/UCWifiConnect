package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
)

func main() {
	c := netman.DefaultClient()
	ssidsFile := os.Getenv("SNAP_COMMON") + "/ssids"
	for {
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
			}
		}
		// wait 5 seconds
		time.Sleep(5000 * time.Millisecond)
	}
}
