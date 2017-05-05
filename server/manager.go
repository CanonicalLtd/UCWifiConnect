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

	"github.com/CanonicalLtd/UCWifiConnect/utils"
)

const (
	address = ":8080"
)

const (
	// NONE any server is up
	NONE RunningServer = 0 + iota
	// MANAGEMENT only management portal is up. Rest are down
	MANAGEMENT
	// OPERATIONAL only operational portal is up. Rest are down
	OPERATIONAL
)

// RunningServer enum defining which server is up and running
type RunningServer int

var currentlyRunning = NONE
var managementCloser io.Closer
var operationalCloser io.Closer

// Running returns RunningServer enum value saying which server is running
func Running() RunningServer {
	return currentlyRunning
}

// StartManagementServer starts server in management mode
func StartManagementServer() error {
	utils.WriteWaitFile()
	var err error
	managementCloser, err = listenAndServe(address, managementHandler())
	if err != nil {
		managementCloser = nil
		return err
	}
	currentlyRunning = MANAGEMENT
	return nil
}

// StartOperationalServer starts server in operational mode
func StartOperationalServer() error {
	var err error
	operationalCloser, err = listenAndServe(address, operationalHandler())
	if err != nil {
		operationalCloser = nil
		return err
	}
	currentlyRunning = OPERATIONAL
	return nil
}

// ShutdownManagementServer shutdown server management mode. If server is in operational mode, it does nothing
func ShutdownManagementServer() error {
	if managementCloser == nil {
		log.Print("== wifi-connect: Skipping stop management server since it is not up")
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
	currentlyRunning = NONE
	return nil
}

// ShutdownOperationalServer shutdown server operational mode. If server is up in management mode, it does nothing
func ShutdownOperationalServer() error {
	if operationalCloser == nil {
		log.Print("== wifi-connect: Skipping stop operational server since it is not up")
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
	currentlyRunning = NONE
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
	currentlyRunning = NONE
	return nil
}
