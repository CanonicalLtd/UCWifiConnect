package main

import (
	//"bufio"
	//"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
)

/*
type options struct {
	devel bool
}

func args() *options {
	opts := &options{}
	flag.BoolVar(&opts.devel, "devel", false, "Test a hard coded devel path")
	return opts
}
*/
func main() {
	c := netman.DefaultClient()
	/*
		opts := args()
			if opts.devel {
				dvs := c.GetDevices()
				dwvs := c.GetWifiDevices(dvs)
				for _, d := range dwvs {
					fmt.Println(d)
				}
				return
			}
	*/
	ssidsFile := os.Getenv("SNAP_COMMON") + "/ssids"
	for {
		fmt.Println("==== iter")
		time.Sleep(200 * time.Millisecond)
		SSIDs, _, _ := c.Ssids()
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
	return
}
