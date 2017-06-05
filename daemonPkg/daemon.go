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

package deamonPkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/CanonicalLtd/UCWifiConnect/avahi"
	"github.com/CanonicalLtd/UCWifiConnect/netman"
	"github.com/CanonicalLtd/UCWifiConnect/server"
	"github.com/CanonicalLtd/UCWifiConnect/utils"
	"github.com/CanonicalLtd/UCWifiConnect/wifiap"
)

// enum to track current system State
const (
	STARTING = 0 + iota
	MANAGING
	OPERATING
	MANUAL
)

var ManualFlagPath string
var WaitFlagPath string
var PreviousState = STARTING
var State = STARTING

//used to clase the operataional http server
var err error

func SetState(s int) {
	PreviousState = State
	State = s
}

// ScanSsids sets wlan0 to be managed and then scans
// for ssids. If found, write the ssids (comma separated)
// to path and return true, else return false.
func ScanSsids(path string, c *netman.Client) bool {
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

// Unmanage sets wlan0 to be Unmanaged by network manager if it
// is managed
func Unmanage(c *netman.Client) {
	ifaces, _ := c.WifisManaged(c.GetWifiDevices(c.GetDevices()))
	if _, ok := ifaces["wlan0"]; ok {
		c.SetIfaceManaged("wlan0", false, c.GetWifiDevices(c.GetDevices()))
	}
}

// manage sets wlan0 to not managed by network manager
func manage(c *netman.Client) {
	c.SetIfaceManaged("wlan0", true, c.GetWifiDevices(c.GetDevices()))
}

// CheckWaitApConnect returns true if the flag wait file exists
// and false if it does not
func CheckWaitApConnect() bool {
	if _, err := os.Stat(WaitFlagPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// CheckManualMode returns true if the manual mode flag wait file exists
// and false if it does not
func CheckManualMode() bool {
	if _, err := os.Stat(ManualFlagPath); os.IsNotExist(err) {
		if State == MANUAL {
			SetState(STARTING)
			fmt.Println("== wifi-connect: entering STARTING mode")
		}
		return false
	}
	if State != MANUAL {
		SetState(MANUAL)
		fmt.Println("== wifi-connect: entering MANUAL mode")
	}
	return true
}

// if wifiap is UP and there are no known SSIDs, bring it down so on next
// loop iter we start again and can get SSIDs. returns true when ip is
// UP and has no ssids
func IsApUpWithoutSSIDs(cw *wifiap.Client) bool {
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

// ManagementServerUp starts the management server if it is
// not running
func ManagementServerUp() {
	if server.Current != server.Management && server.State == server.Stopped {
		err = server.StartManagementServer()
		if err != nil {
			fmt.Println("== wifi-connect: Error start Mamagement portal:", err)
		}
		// init mDNS
		avahi.InitMDNS()
	}
}

// ManagementServerDown stops the management server if it is running
// also remove the wait flag file, thus resetting proper State
func ManagementServerDown() {
	if server.Current == server.Management && (server.State == server.Running || server.State == server.Starting) {
		err = server.ShutdownManagementServer()
		if err != nil {
			fmt.Println("== wifi-connect: Error stopping the Management portal:", err)
		}
		//remove flag fie so daemon resumes normal control
		utils.RemoveFlagFile(os.Getenv("SNAP_COMMON") + "/startingApConnect")
	}
}

// OperationalServerUp starts the operational server if it is
// not running
func OperationalServerUp() {
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
func OperationalServerDown() {
	if server.Current == server.Operational && (server.State == server.Running || server.State == server.Starting) {
		err = server.ShutdownOperationalServer()
		if err != nil {
			fmt.Println("== wifi-connect: Error stopping Operational portal:", err)
		}
	}
}

// SetDefaults sets defaults if not yet set. Currently the hash
// for the portals password is set.
// TODO: set default password based on MAC addr or Serial number
func SetDefaults() {
	if _, err := os.Stat(utils.HashFile); os.IsNotExist(err) {
		utils.HashIt("wifi-connect")
	}
}
