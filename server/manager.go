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

package server

import (
	"os"

	"fmt"

	"github.com/CanonicalLtd/UCWifiConnect/utils"
)

// Enum of available server options
const (
	None Server = 0 + iota
	Management
	Operational
)

// Server defines an enum of servers
type Server int

// Current active server instance. None if any is enabled at this moment
var Current = None

// StartManagementServer starts server in management mode
func StartManagementServer() error {
	if Current != None {
		Current = None
		return fmt.Errorf("Not in a valid status. Please stop first any other server instance before starting this one")
	}

	// change current instance asap we manage this server
	Current = Management

	waitPath := os.Getenv("SNAP_COMMON") + "/startingApConnect"
	err := utils.WriteFlagFile(waitPath)
	if err != nil {
		Current = None
		return err
	}

	err = listenAndServe(address, managementHandler())
	if err != nil {
		Current = None
		return err
	}

	return nil
}

// StartOperationalServer starts server in operational mode
func StartOperationalServer() error {
	if Current != None {
		return fmt.Errorf("Not in a valid status. Please stop first any other server instance before starting this one")
	}

	// change current instance asap we manage this server
	Current = Operational

	err := listenAndServe(address, operationalHandler())
	if err != nil {
		Current = None
		return err
	}

	return nil
}

// ShutdownManagementServer shutdown server management mode. If management server is not up, returns error
func ShutdownManagementServer() error {
	if Current != Management || (State != Running && State != Starting) {
		return fmt.Errorf("Trying to stop management server when it is not running")
	}

	err := stop()
	if err != nil {
		return err
	}

	Current = None
	return nil
}

// ShutdownOperationalServer shutdown server operational mode. If operational server is not up, returns error
func ShutdownOperationalServer() error {
	if Current != Operational || (State != Running && State != Starting) {
		return fmt.Errorf("Trying to stop operational server when it is not running")
	}

	err := stop()
	if err != nil {
		return err
	}

	Current = None
	return nil
}
