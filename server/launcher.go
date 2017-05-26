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
	"log"
	"net"
	"net/http"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/utils"
)

const (
	address = ":8080"
)

// Server running state
const (
	Stopped RunningState = 0 + iota
	Starting
	Running
	Stopping
)

// RunningState enum defining which server is up and running
type RunningState int

// State holds current server state
var State = Stopped

var listener net.Listener
var done chan bool

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// WaitForState waits for server reach certain state
func WaitForState(state RunningState) bool {
	retries := 10
	idle := 10 * time.Millisecond
	for ; retries > 0 && State != state; retries-- {
		time.Sleep(idle)
		idle *= 2
	}
	return State == state
}

// Accept accepts incoming tcp connections
func (ln tcpKeepAliveListener) Accept() (net.Conn, error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return tc, err
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func listenAndServe(addr string, handler http.Handler) error {

	if State != Stopped {
		return Errorf("Server is not in proper stopped state before trying to start it")
	}

	if utils.RunningOn(addr) {
		return Errorf("Another instance is running in same address %v", addr)
	}

	State = Starting

	srv := &http.Server{Addr: addr, Handler: handler}
	// channel needed to communicate real server shutdown, as after calling listener.Close()
	// it can take several milliseconds to really stop the listening.
	done = make(chan bool)

	var err error
	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// launch goroutine to check server state changes after startup is triggered
	go func() {
		retries := 10
		idle := 10 * time.Millisecond
		for ; !utils.RunningOn(addr) && retries > 0; retries-- {
			time.Sleep(idle)
			idle *= 2
		}

		if retries == 0 {
			log.Print(Sprintf("Server could not be started"))
			return
		}

		State = Running
	}()

	// launching server in a goroutine for not blocking
	go func() {
		if listener != nil {
			err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
			if err != nil {
				log.Printf(Sprintf("HTTP Server closing - %v", err))
			}
			// notify server real stop
			done <- true
		}

		close(done)
	}()

	return nil
}

func stop() error {

	if State == Stopped {
		return Errorf("Already stopped")
	}

	if listener == nil {
		State = Stopped
		return Errorf("Already closed")
	}

	State = Stopping

	err := listener.Close()
	if err != nil {
		return err
	}
	listener = nil

	// wait for server real shutdown confirmation
	<-done

	State = Stopped
	return nil
}
