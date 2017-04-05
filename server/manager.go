package server

import (
	"io"
	"log"
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
	// ALL operational and management portals are both up.
	ALL
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
	var err error
	managementCloser, err = listenAndServe(address, managementHandler())
	if err != nil {
		managementCloser = nil
		return err
	}
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
	return nil
}

// ShutdownManagementServer shutdown server management mode. If server is in operational mode, it does nothing
func ShutdownManagementServer() error {
	if managementCloser == nil {
		log.Print("Skipping stop management server since it is not up")
		return nil
	}

	err := managementCloser.Close()
	managementCloser = nil
	return err
}

// ShutdownOperationalServer shutdown server operational mode. If server is up in management mode, it does nothing
func ShutdownOperationalServer() error {
	if operationalCloser == nil {
		log.Print("Skipping stop operational server since it is not up")
		return nil
	}

	err := operationalCloser.Close()
	operationalCloser = nil
	return err
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
	return nil
}
