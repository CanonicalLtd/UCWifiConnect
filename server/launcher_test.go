package server

import (
	"testing"

	telnet "github.com/reiver/go-telnet"
)

func TestLaunchAndStop(t *testing.T) {

	thePort := ":14444"

	srv, err := listenAndServe(thePort, nil)
	if err != nil {
		t.Errorf("Start server failed: %v", err)
	}

	if srv == nil {
		t.Error("Server could not be initialzed")
	}

	// telnet to check server is alive
	caller := telnet.StandardCaller
	err = telnet.DialToAndCall("localhost"+thePort, caller)
	if err != nil {
		t.Errorf("Failed to telnet localhost server at port %v: %v", thePort, err)
	}

	err = srv.Close()
	if err != nil {
		t.Errorf("Stop server error: %v", err)
	}

}
