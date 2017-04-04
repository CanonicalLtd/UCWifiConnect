package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
	"github.com/CanonicalLtd/UCWifiConnect/server"
	"github.com/CanonicalLtd/UCWifiConnect/wifiap"
)

// enum to track current system state
const (
	unknown = iota
	managed
	operational
)

//used to clase the management http server
var mgmtCloser io.Closer

//used to clase the operataional http server
var operCloser io.Closer
var err error

func scanSsids(path string, c *netman.Client) bool {
	manage(c)
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

func managementPortalUp() (io.Closer, error) {
	mgmtCloser, err = server.StartManagementServer()
	if err != nil {
		return mgmtCloser, err
	}
	return mgmtCloser, nil
}
func operationalPortalUp() (io.Closer, error) {
	operCloser, err = server.StartExternalServer()
	if err != nil {
		return operCloser, err
	}
	return operCloser, nil
}

func portalDown(closer io.Closer) error {
	err = closer.Close()
	if err != nil {
		return err
	}
	return nil
}

func unmanage(c *netman.Client) {
	ifaces, _ := c.WifisManaged(c.GetWifiDevices(c.GetDevices()))
	if _, ok := ifaces["wlan0"]; ok {
		fmt.Println("==== Setting wlan0 unmanaged")
		c.SetIfaceUnmanaged("wlan0", c.GetWifiDevices(c.GetDevices()))
	}
}

func manage(c *netman.Client) {
	ifaces, _ := c.WifisManaged(c.GetWifiDevices(c.GetDevices()))
	if _, ok := ifaces["wlan0"]; ok {
		fmt.Println("==== Setting wlan0 managed")
		c.SetIfaceManaged("wlan0", c.GetWifiDevices(c.GetDevices()))
	}
}

func main() {
	first := true
	c := netman.DefaultClient()
	cw := wifiap.DefaultClient()
	ssidsPath := os.Getenv("SNAP_COMMON") + "/ssids"
	for {
		if first {
			first = false
			//wait time period (TBD) on first run to allow wifi connections
			time.Sleep(5000 * time.Millisecond)
		}
		// wait 5 seconds on each iter
		time.Sleep(60000 * time.Millisecond)
		//if an external wifi connection: Operational mode
		fmt.Println("== wifi connected:", c.ConnectedWifi(c.GetWifiDevices(c.GetDevices())))
		if c.ConnectedWifi(c.GetWifiDevices(c.GetDevices())) {
			// Operational mode
			fmt.Println("==== Operational Mode ")
			fmt.Println("==== Stop Management Mode http server...")
			if mgmtCloser != nil {
				err = portalDown(mgmtCloser)
				if err != nil {
					fmt.Println("Error stopping the Management portal:", err)
				}
			}
			fmt.Println("==== Start Operational Mode http server...")
			operCloser, err = operationalPortalUp()
			if err != nil {
				fmt.Println("Error starting the Operational portal:", err)
			}
			continue
		}

		fmt.Println("==== Managment Mode ")
		//is wlan0 managed? yes, set unmanaged
		unmanage(c)

		//wifi-ap UP?
		wifiUp, err := cw.Enabled()
		if err != nil {
			fmt.Println("==== Error checking wifi-ap.Enabled():", err)
			continue // try again since no better course of action
		}
		fmt.Println("==== wifi-ap enabled state:", wifiUp)
		if !wifiUp {
			found := scanSsids(ssidsPath, c)
			unmanage(c)
			if !found {
				fmt.Println("==== No SSIDs found. Looping...")
				//continue
			}
			fmt.Println("==== Stop Management portal")
			if mgmtCloser != nil {
				err = portalDown(mgmtCloser)
				if err != nil {
					fmt.Println("==== Error stopping Mamagement portal:", err)
					continue // try again since no better course of action
				}
			}
			fmt.Println("==== Start wifi-ap")
			cw.Enable()
		}

		fmt.Println("==== Start Management portal")
		mgmtCloser, err = managementPortalUp()
		if err != nil {
			fmt.Println("==== Error start Mamagement portal:", err)
			continue // try again for lack of better option
		}
	}
}
