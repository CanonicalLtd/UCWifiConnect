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
