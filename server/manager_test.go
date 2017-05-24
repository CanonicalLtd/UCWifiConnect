package server

import (
	"os"
	"testing"
	"time"
)

// needed method for testing, as server state is updated asynchronously
func waitForState(state RunningState) bool {
	retries := 10
	idle := 1000 * time.Millisecond
	for ; retries > 0 && State != state; retries-- {
		time.Sleep(idle)
	}
	return State == state
}

func TestBasicServerTransitionStates(t *testing.T) {

	os.Setenv("SNAP_COMMON", os.TempDir())

	if Current != None || State != Stopped {
		t.Errorf("Server is not in initial state")
	}

	if err := StartManagementServer(); err != nil {
		t.Errorf("Error starting management server %v", err)
	}
	if Current != Management || (State != Starting && State != Running) {
		t.Errorf("Server is not in starting or in management status")
	}

	waitForState(Running)

	if err := ShutdownManagementServer(); err != nil {
		t.Errorf("Error stopping management server %v", err)
	}

	if Current != None {
		t.Errorf("Current server is not None")
	}

	waitForState(Stopped)

	if err := StartOperationalServer(); err != nil {
		t.Errorf("Error starting operational server %v", err)
	}
	if Current != Operational || (State != Starting && State != Running) {
		t.Errorf("Server is not in starting or in operational status")
	}

	waitForState(Running)

	if err := ShutdownOperationalServer(); err != nil {
		t.Errorf("Error stopping operational server %v", err)
	}
	if Current != None {
		t.Errorf("Current server is not none")
	}
}

func TestEdgeServerTransitionStates(t *testing.T) {
	os.Setenv("SNAP_COMMON", os.TempDir())

	if Current != None {
		t.Errorf("Server is not in initial state")
	}

	if err := StartManagementServer(); err != nil {
		t.Errorf("Error starting management server %v", err)
	}
	if Current != Management || (State != Starting && State != Running) {
		t.Errorf("Server is not in starting or in management status")
	}

	waitForState(Running)

	// start operational server without stopping management must throw an error
	if err := StartOperationalServer; err == nil {
		t.Errorf(`Expected an error when trying to launch one server instance having 
		the other active`)
	}
	if Current != Management {
		t.Errorf("Server is not in management status after failed start operational server")
	}

	// stop wrong server must throw an error
	if err := ShutdownOperationalServer; err == nil {
		t.Errorf("Expected an error when trying to shutdown wrong server")
	}
	if Current != Management {
		t.Errorf("Server is not in management status after failed start operational server")
	}

	if err := ShutdownManagementServer(); err != nil {
		t.Errorf("Error stopping management server %v", err)
	}
	if Current != None {
		t.Errorf("Server is not in None status")
	}

	waitForState(Stopped)

	// analog tests with operational server
	if err := StartOperationalServer(); err != nil {
		t.Errorf("Error starting operational server %v", err)
	}
	if Current != Operational || (State != Starting && State != Running) {
		t.Errorf("Server is not in starting or in operational status")
	}

	waitForState(Running)

	// start management server without stopping operational must throw an error
	if err := StartManagementServer; err == nil {
		t.Errorf(`Expected an error when trying to launch one server instance having 
		the other active`)
	}
	if Current != Operational {
		t.Errorf("Server is not in operational status after failed start operational server")
	}

	// stop wrong server must throw an error
	if err := ShutdownManagementServer; err == nil {
		t.Errorf("Expected an error when trying to shutdown wrong server")
	}
	if Current != Operational {
		t.Errorf("Server is not in operational status after failed start operational server")
	}

	if err := ShutdownOperationalServer(); err != nil {
		t.Errorf("Error stopping operational server %v", err)
	}
	if Current != None {
		t.Errorf("Server is not in None status")
	}
}
