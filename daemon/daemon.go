/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/avahi"
	"github.com/CanonicalLtd/UCWifiConnect/netman"
	"github.com/CanonicalLtd/UCWifiConnect/server"
	"github.com/CanonicalLtd/UCWifiConnect/utils"
	"github.com/CanonicalLtd/UCWifiConnect/wifiap"
)

// enum to track current system state
const (
	unknown = iota
	managed
	operational
)

//used to clase the operataional http server
var err error

// scanSsids sets wlan0 to be managed and then scans
// for ssids. If found, write the ssids (comma separated)
// to path and return true, else return false.
func scanSsids(path string, c *netman.Client) bool {
	manage(c)
	time.Sleep(5000 * time.Millisecond)
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
			fmt.Println("==== Error writing SSID(s) to ", path)
		} else {
			fmt.Println("==== SSID(s) found and written to ", path)
			return true
		}
	}
	fmt.Println("==== NO SSID found")
	return false
}

// unmanage sets wlan0 to be unmanaged by network manager if it
// is managed
func unmanage(c *netman.Client) {
	ifaces, _ := c.WifisManaged(c.GetWifiDevices(c.GetDevices()))
	if _, ok := ifaces["wlan0"]; ok {
		fmt.Println("==== Setting wlan0 unmanaged")
		c.SetIfaceManaged("wlan0", false, c.GetWifiDevices(c.GetDevices()))
	}
}

// manage sets wlan0 to not managed by network manager
func manage(c *netman.Client) {
	fmt.Println("==== Setting wlan0 managed")
	c.SetIfaceManaged("wlan0", true, c.GetWifiDevices(c.GetDevices()))
}

// checkWaitApConnect returns true if the flag wait file exists
// and false if it does not
func checkWaitApConnect() bool {
	waitApPath := os.Getenv("SNAP_COMMON") + "/startingApConnect"
	if _, err := os.Stat(waitApPath); os.IsNotExist(err) {
		fmt.Println("==== Wait file not found")
		return false
	}
	fmt.Println("==== Wait file found")
	return true
}

// managementServerUp starts the management server if it is
// not running
func managementServerUp() {
	if server.Running() != server.MANAGEMENT {
		err = server.StartManagementServer()
		if err != nil {
			fmt.Println("==== Error start Mamagement portal:", err)
		}
	}
}

// managementServerDown stops the management server if it is running
// also remove the wait flag file, thus resetting proper state
func managementServerDown() {
	if server.Running() == server.MANAGEMENT {
		err = server.ShutdownManagementServer()
		if err != nil {
			fmt.Println("==== Error stopping the Management portal:", err)
		}
		//remove flag fie so daemon resumes normal control
		utils.RemoveWaitFile()
	}
}

// operationalServerUp starts the operational server if it is
// not running
func operationalServerUp() {
	if server.Running() != server.OPERATIONAL {
		err = server.StartOperationalServer()
		if err != nil {
			fmt.Println("==== Error starting the Operational portal:", err)
		}
	}
}

// operationalServerdown stops the operational server if it is running
func operationalServerDown() {
	if server.Running() == server.OPERATIONAL {
		err = server.ShutdownOperationalServer()
		if err != nil {
			fmt.Println("==== Error stopping Operational portal:", err)
		}
	}
}

func main() {
	first := true
	c := netman.DefaultClient()
	cw := wifiap.DefaultClient()
	ssidsPath := os.Getenv("SNAP_COMMON") + "/ssids"

	// stop servers if running at start
	managementServerDown()
	operationalServerDown()

	//remove possibly left over wait flag file
	utils.RemoveWaitFile()

	for {
		if first {
			first = false
			//wait time period (TBD) on first run to allow wifi connections
			time.Sleep(40000 * time.Millisecond)
			//remove previous state flag, if any on deamon startup
			utils.RemoveWaitFile()
		}

		// wait 5 seconds on each iter
		time.Sleep(5000 * time.Millisecond)

		// wait/loop until management portal's wait flag file is gone
		// this stops daemon state changing until the management portal
		// is done, either stopped or the user has attempted to connect to
		// an external AP
		if checkWaitApConnect() {
			continue
		}

		// if an external wifi connection, we are in Operational mode
		// and we stay here until there is an external wifi connection
		if c.ConnectedWifi(c.GetWifiDevices(c.GetDevices())) {
			fmt.Println("======== Operational Mode ")
			fmt.Println("==== Stop Management Mode http server if running")
			managementServerDown()
			//comment out operational server as later phase of work
			//fmt.Println("==== Start Operational Mode http server if not running")
			//operationalServerUp()
			//avahi.InitMDNS()
			continue
		}

		fmt.Println("====== Management Mode")
		// if wlan0 managed, set unmanaged so that we can bring up wifi-ap
		// properly
		unmanage(c)

		// stop operational portal
		//comment out operational server as later phase of work
		//fmt.Println("==== Stop Operational Mode http server if running")
		//operationalServerDown()

		//wifi-ap UP?
		wifiUp, err := cw.Enabled()
		if err != nil {
			fmt.Println("==== Error checking wifi-ap.Enabled():", err)
			continue // try again since no better course of action
		}

		fmt.Println("==== Wifi-ap enabled state:", wifiUp)

		//get ssids if wifi-ap Down
		if !wifiUp {
			found := scanSsids(ssidsPath, c)
			unmanage(c)
			if !found {
				fmt.Println("==== No SSIDs found. Looping.")
				continue
			}
			fmt.Println("==== Have SSIDs: start wifi-ap")
			cw.Enable()
		}

		fmt.Println("==== Start Management portal if not running")
		managementServerUp()
		//TODO POC avahi
		avahi.InitMDNS()
	}
}
