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
	"os"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/daemon"
	"github.com/CanonicalLtd/UCWifiConnect/netman"
	"github.com/CanonicalLtd/UCWifiConnect/utils"
	"github.com/CanonicalLtd/UCWifiConnect/wifiap"
)

func main() {

	client := daemon.GetClient()
	client.SetDefaults()
	first := true
	client.SetWaitFlagPath(os.Getenv("SNAP_COMMON") + "/startingApConnect")
	client.SetManualFlagPath(os.Getenv("SNAP_COMMON") + "/manualMode")

	c := netman.DefaultClient()
	cw := wifiap.DefaultClient()

	client.ManagementServerDown()
	client.OperationalServerDown()

	for {
		if first {
			fmt.Println("== wifi-connect: daemon STARTING")
			client.SetPreviousState(daemon.STARTING)
			client.SetState(daemon.STARTING)
			first = false
			//clean start require wifi AP down so we can get SSIDs
			cw.Disable()
			//remove previous State flags
			utils.RemoveFlagFile(client.GetWaitFlagPath())
			utils.RemoveFlagFile(client.GetManualFlagPath())
			//TODO only wait if wlan0 is managed
			//wait time period (TBD) on first run to allow wifi connections
			time.Sleep(40000 * time.Millisecond)
		}

		// wait 5 seconds on each iter
		time.Sleep(5000 * time.Millisecond)

		// loop without action if in manual mode
		if client.ManualMode() {
			continue
		}

		// start clean on exiting manual mode
		if client.GetPreviousState() == daemon.MANUAL {
			first = true
			continue
		}
		// the AP should not be up without SSIDS
		if client.IsApUpWithoutSSIDs(cw) {
			cw.Disable()
			continue
		}

		// wait/loop until management portal's wait flag file is gone
		// this stops daemon State changing until the management portal
		// is done, either stopped or the user has attempted to connect to
		// an external AP
		if client.CheckWaitApConnect() {
			continue
		}

		// if an external wifi connection, we are in Operational mode
		// and we stay here until there is an external wifi connection
		if c.ConnectedWifi(c.GetWifiDevices(c.GetDevices())) {
			client.SetState(daemon.OPERATING)
			if client.GetPreviousState() != daemon.OPERATING {
				fmt.Println("== wifi-connect: entering OPERATIONAL mode")
			}
			if client.GetPreviousState() == daemon.MANAGING {
				client.ManagementServerDown()
			}
			client.OperationalServerUp()
			continue
		}

		client.SetState(daemon.MANAGING)
		if client.GetPreviousState() != daemon.MANAGING {
			fmt.Println("== wifi-connect: entering MANAGEMENT mode")
		}

		// if wlan0 managed, set Unmanaged so that we can bring up wifi-ap
		// properly
		client.Unmanage(c)

		//wifi-ap UP?
		wifiUp, err := cw.Enabled()
		if err != nil {
			fmt.Println("== wifi-connect: Error checking wifi-ap.Enabled():", err)
			continue // try again since no better course of action
		}

		//get ssids if wifi-ap Down
		if !wifiUp {
			found := client.ScanSsids(utils.SsidsFile, c)
			client.Unmanage(c)
			if !found {
				fmt.Println("== wifi-connect: Looping.")
				continue
			}
			fmt.Println("== wifi-connect: starting wifi-ap")
			cw.Enable()
			if client.GetPreviousState() == daemon.OPERATING {
				client.OperationalServerDown()
			}
			client.ManagementServerUp()
		}
	}
}
