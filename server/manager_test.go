package server

import (
	"os"
	"testing"
	"time"
)

// needed method for testing, as server state is updated asynchronously
func waitForState(state State) bool {
	retries := 10
	idle := 1000 * time.Millisecond
	for ; retries > 0 && currentlyRunning != state; retries-- {
		time.Sleep(idle)
	}
	return currentlyRunning == state
}

func TestBasicServerTransitionStates(t *testing.T) {

	os.Setenv("SNAP_COMMON", os.TempDir())

	if Running() != None {
		t.Errorf("Server is not in initial state")
	}

	if err := StartManagementServer(); err != nil {
		t.Errorf("Error starting management server %v", err)
	}
	if Running() != Management && Running() != StartingManagement {
		t.Errorf("Server is not in starting or in management status")
	}

	waitForState(Management)

	if err := ShutdownManagementServer(); err != nil {
		t.Errorf("Error stopping management server %v", err)
	}

	if Running() != None {
		t.Errorf("Server is not in None status")
	}

	waitForState(None)

	time.Sleep(5 * time.Second)

	if err := StartOperationalServer(); err != nil {
		t.Errorf("Error starting operational server %v", err)
	}
	if Running() != Operational && Running() != StartingOperational {
		t.Errorf("Server is not in starting or in operational status")
	}

	time.Sleep(5 * time.Second)

	//waitForState(Operational)

	if err := ShutdownOperationalServer(); err != nil {
		t.Errorf("Error stopping operational server %v", err)
	}
	if Running() != None {
		t.Errorf("Server is not in None status")
	}
}

func TestEdgeServerTransitionStates(t *testing.T) {
	os.Setenv("SNAP_COMMON", os.TempDir())

	if Running() != None {
		t.Errorf("Server is not in initial state")
	}

	if err := StartManagementServer(); err != nil {
		t.Errorf("Error starting management server %v", err)
	}
	if Running() != Management && Running() != StartingManagement {
		t.Errorf("Server is not in starting or in management status")
	}

	// start operational server without stopping management must throw an error
	if err := StartOperationalServer; err == nil {
		t.Errorf(`Expected an error when trying to launch one server instance having 
		the other active`)
	}
	if Running() != Management {
		t.Errorf("Server is not in management status after failed start operational server")
	}

	// stop wrong server must throw an error
	if err := ShutdownOperationalServer; err == nil {
		t.Errorf("Expected an error when trying to shutdown wrong server")
	}
	if Running() != Management {
		t.Errorf("Server is not in management status after failed start operational server")
	}

	if err := ShutdownManagementServer(); err != nil {
		t.Errorf("Error stopping management server %v", err)
	}
	if Running() != None {
		t.Errorf("Server is not in None status")
	}

	// analog tests with operational server
	if err := StartOperationalServer(); err != nil {
		t.Errorf("Error starting operational server %v", err)
	}
	if Running() != Operational && Running() != StartingOperational {
		t.Errorf("Server is not in starting or in operational status")
	}

	// start management server without stopping operational must throw an error
	if err := StartManagementServer; err == nil {
		t.Errorf(`Expected an error when trying to launch one server instance having 
		the other active`)
	}
	if Running() != Operational {
		t.Errorf("Server is not in operational status after failed start operational server")
	}

	// stop wrong server must throw an error
	if err := ShutdownManagementServer; err == nil {
		t.Errorf("Expected an error when trying to shutdown wrong server")
	}
	if Running() != Operational {
		t.Errorf("Server is not in operational status after failed start operational server")
	}

	if err := ShutdownOperationalServer(); err != nil {
		t.Errorf("Error stopping operational server %v", err)
	}
	if Running() != None {
		t.Errorf("Server is not in None status")
	}
}