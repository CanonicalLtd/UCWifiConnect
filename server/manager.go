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
	"io"
	"log"
	"os"

	"time"

	"fmt"

	"github.com/CanonicalLtd/UCWifiConnect/utils"
)

const (
	address = ":8080"
)

// Server states
const (
	None State = 0 + iota
	StartingManagement
	ShuttingDownManagement
	Management
	StartingOperational
	ShuttingDownOperational
	Operational
)

// State enum defining which server is up and running
type State int

var currentlyRunning = None
var managementCloser io.Closer
var operationalCloser io.Closer

// Running returns ServerState enum value saying which server is running
func Running() State {
	return currentlyRunning
}

func updateStateWhenServerUp() {
	go func() {
		var i int
		for i = 0; !utils.RunningOn(address) && i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
		}

		if i < 0 {
			log.Print("Server could not be started")
			return
		}

		switch currentlyRunning {
		case StartingManagement:
			currentlyRunning = Management
		case StartingOperational:
			currentlyRunning = Operational
		case ShuttingDownManagement:
			fallthrough
		case ShuttingDownOperational:
			currentlyRunning = None
		}
	}()
}

// StartManagementServer starts server in management mode
func StartManagementServer() error {
	if currentlyRunning != None {
		return fmt.Errorf("Not in a valid status. Please stop first any other server instance before starting this one")
	}

	currentlyRunning = StartingManagement

	waitPath := os.Getenv("SNAP_COMMON") + "/startingApConnect"
	var err error
	err = utils.WriteFlagFile(waitPath)
	if err != nil {
		return err
	}

	managementCloser, err = listenAndServe(address, managementHandler())
	if err != nil {
		managementCloser = nil
		return err
	}

	updateStateWhenServerUp()
	return nil
}

// StartOperationalServer starts server in operational mode
func StartOperationalServer() error {
	if currentlyRunning != None {
		return fmt.Errorf("Not in a valid status. Please stop first any other server instance before starting this one")
	}

	currentlyRunning = StartingOperational

	var err error
	operationalCloser, err = listenAndServe(address, operationalHandler())
	if err != nil {
		operationalCloser = nil
		return err
	}

	updateStateWhenServerUp()
	return nil
}

// ShutdownManagementServer shutdown server management mode. If management server is not up, returns error
func ShutdownManagementServer() error {
	if currentlyRunning != Management && currentlyRunning != StartingManagement {
		return fmt.Errorf("Trying to stop management server when it is not running")
	}

	if managementCloser == nil {
		return nil
	}

	err := managementCloser.Close()
	if err != nil {
		return err
	}
	managementCloser = nil
	// TODO for now we only have one server up at a time. Later, if happens
	// that more than one can be up at the same time it would be needed manage this
	// state changes in a better way
	currentlyRunning = None
	return nil
}

// ShutdownOperationalServer shutdown server operational mode. If operational server is not up, returns error
func ShutdownOperationalServer() error {
	if currentlyRunning != Operational && currentlyRunning != StartingOperational {
		return fmt.Errorf("Trying to stop operational server when it is not running")
	}

	if operationalCloser == nil {
		return nil
	}

	err := operationalCloser.Close()
	if err != nil {
		return err
	}
	operationalCloser = nil
	// TODO for now we only have one server up at a time. Later, if happens
	// that more than one can be up at the same time it would be needed manage this
	// state changes in a better way
	currentlyRunning = None
	return nil
}

// ShutdownServer shutdown server no matter the mode it is up
func ShutdownServer() error {
	err := ShutdownManagementServer()
	err2 := ShutdownOperationalServer()
	if err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}
	currentlyRunning = None
	return nil
}
