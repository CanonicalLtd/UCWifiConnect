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
	STARTING = 0 + iota
	MANAGING
	OPERATING
	MANUAL
)

var manualFlagPath string
var waitFlagPath string
var previousState = STARTING
var state = STARTING

//used to clase the operataional http server
var err error

func setState(s int) {
	previousState = state
	state = s
}

// scanSsids sets wlan0 to be managed and then scans
// for ssids. If found, write the ssids (comma separated)
// to path and return true, else return false.
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
			fmt.Println("== wifi-connect: Error writing SSID(s) to ", path)
		} else {
			fmt.Println("== wifi-connect: SSID(s) obtained")
			return true
		}
	}
	fmt.Println("== wifi-connect: No SSID found")
	return false
}

// unmanage sets wlan0 to be unmanaged by network manager if it
// is managed
func unmanage(c *netman.Client) {
	ifaces, _ := c.WifisManaged(c.GetWifiDevices(c.GetDevices()))
	if _, ok := ifaces["wlan0"]; ok {
		c.SetIfaceManaged("wlan0", false, c.GetWifiDevices(c.GetDevices()))
	}
}

// manage sets wlan0 to not managed by network manager
func manage(c *netman.Client) {
	c.SetIfaceManaged("wlan0", true, c.GetWifiDevices(c.GetDevices()))
}

// checkWaitApConnect returns true if the flag wait file exists
// and false if it does not
func checkWaitApConnect() bool {
	if _, err := os.Stat(waitFlagPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// checkManualMode returns true if the manual mode flag wait file exists
// and false if it does not
func checkManualMode() bool {
	if _, err := os.Stat(manualFlagPath); os.IsNotExist(err) {
		if state == MANUAL {
			setState(STARTING)
			fmt.Println("== wifi-connect: entering STARTING mode")
		}
		return false
	}
	if state != MANUAL {
		setState(MANUAL)
		fmt.Println("== wifi-connect: entering MANUAL mode")
	}
	return true
}

// if wifiap is UP and there are no known SSIDs, bring it down so on next
// loop iter we start again and can get SSIDs. returns true when ip is
// UP and has no ssids
func isApUpWithoutSSIDs(cw *wifiap.Client) bool {
	wifiUp, _ := cw.Enabled()
	if !wifiUp {
		return false
	}
	ssids, _ := utils.ReadSsidsFile()
	if len(ssids) < 1 {
		fmt.Println("== wifi-connect: wifi-ap is UP but has no SSIDS")
		return true // ap is up with no ssids
	}
	return false
}

// managementServerUp starts the management server if it is
// not running
func managementServerUp() {
	if server.Current != server.Management && server.State == server.Stopped {
		err = server.StartManagementServer()
		if err != nil {
			fmt.Println("== wifi-connect: Error start Mamagement portal:", err)
		}
		// init mDNS
		avahi.InitMDNS()
	}
}

// managementServerDown stops the management server if it is running
// also remove the wait flag file, thus resetting proper state
func managementServerDown() {
	if server.Current == server.Management && (server.State == server.Running || server.State == server.Starting) {
		err = server.ShutdownManagementServer()
		if err != nil {
			fmt.Println("== wifi-connect: Error stopping the Management portal:", err)
		}
		//remove flag fie so daemon resumes normal control
		utils.RemoveFlagFile(os.Getenv("SNAP_COMMON") + "/startingApConnect")
	}
}

// operationalServerUp starts the operational server if it is
// not running
func operationalServerUp() {
	if server.Current != server.Operational && server.State == server.Stopped {
		err = server.StartOperationalServer()
		if err != nil {
			fmt.Println("== wifi-connect: Error starting the Operational portal:", err)
		}
		// init mDNS
		avahi.InitMDNS()
	}
}

// operationalServerdown stops the operational server if it is running
func operationalServerDown() {
	if server.Current == server.Operational && (server.State == server.Running || server.State == server.Starting) {
		err = server.ShutdownOperationalServer()
		if err != nil {
			fmt.Println("== wifi-connect: Error stopping Operational portal:", err)
		}
	}
}

func main() {
	first := true
	waitFlagPath = os.Getenv("SNAP_COMMON") + "/startingApConnect"
	manualFlagPath = os.Getenv("SNAP_COMMON") + "/manualMode"

	c := netman.DefaultClient()
	cw := wifiap.DefaultClient()

	// stop servers if running at start
	managementServerDown()
	operationalServerDown()

	for {
		if first {
			fmt.Println("== wifi-connect: daemon STARTING")
			previousState = STARTING
			state = STARTING
			first = false
			//clean start require wifi AP down so we can get SSIDs
			cw.Disable()
			//TODO only wait if wlan0 is managed
			//remove previous state flag, if any on deamon startup
			utils.RemoveFlagFile(waitFlagPath)
			utils.RemoveFlagFile(manualFlagPath)
			//wait time period (TBD) on first run to allow wifi connections
			time.Sleep(40000 * time.Millisecond)
		}

		// wait 5 seconds on each iter
		time.Sleep(5000 * time.Millisecond)

		// loop without action if in manual mode
		if checkManualMode() {
			continue
		}

		// start clean on exiting manual mode
		if previousState == MANUAL {
			first = true
			continue
		}
		// the AP should not be up without SSIDS
		if isApUpWithoutSSIDs(cw) {
			cw.Disable()
			continue
		}

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
			setState(OPERATING)
			if previousState != OPERATING {
				fmt.Println("== wifi-connect: entering OPERATIONAL mode")
			}
			if previousState == MANAGING {
				managementServerDown()
			}
			operationalServerUp()
			continue
		}

		setState(MANAGING)
		if previousState != MANAGING {
			fmt.Println("== wifi-connect: entering MANAGEMENT mode")
		}

		// if wlan0 managed, set unmanaged so that we can bring up wifi-ap
		// properly
		unmanage(c)

		//wifi-ap UP?
		wifiUp, err := cw.Enabled()
		if err != nil {
			fmt.Println("== wifi-connect: Error checking wifi-ap.Enabled():", err)
			continue // try again since no better course of action
		}

		//get ssids if wifi-ap Down
		if !wifiUp {
			found := scanSsids(utils.SsidsFile, c)
			unmanage(c)
			if !found {
				fmt.Println("== wifi-connect: Looping.")
				continue
			}
			fmt.Println("== wifi-connect: starting wifi-ap")
			cw.Enable()
			if previousState == OPERATING {
				operationalServerDown()
			}
			managementServerUp()
		}
	}
}
