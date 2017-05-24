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
	"testing"

	telnet "github.com/reiver/go-telnet"
)

func TestLaunchAndStop(t *testing.T) {

	thePort := ":14444"

	err := listenAndServe(thePort, nil)
	if err != nil {
		t.Errorf("Start server failed: %v", err)
	}

	// telnet to check server is alive
	caller := telnet.StandardCaller
	err = telnet.DialToAndCall("localhost"+thePort, caller)
	if err != nil {
		t.Errorf("Failed to telnet localhost server at port %v: %v", thePort, err)
	}

	err = stop()
	if err != nil {
		t.Errorf("Stop server error: %v", err)
	}
}

func TestStates(t *testing.T) {

	WaitForState(Stopped)

	if State != Stopped {
		t.Error("Not in initial state")
	}

	thePort := ":14444"

	err := listenAndServe(thePort, nil)
	if err != nil {
		t.Errorf("Start server failed: %v", err)
	}

	if State != Starting && State != Running {
		t.Error("Not in proper start(ing) state")
	}

	WaitForState(Running)

	// try a bad transition
	err = listenAndServe(thePort, nil)
	if err == nil {
		t.Error("An error should be thrown when trying to start an already running instance")
	}

	err = stop()
	if err != nil {
		t.Errorf("Stop server error: %v", err)
	}

	if State != Stopping && State != Stopped {
		t.Error("Not in proper stop(ing) state")
	}

	WaitForState(Stopped)

	// try bad transitions
	err = stop()
	if err == nil {
		t.Error("An error should be thrown when trying to stop a stopped instance")
	}
}
