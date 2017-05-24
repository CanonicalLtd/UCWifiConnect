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
	"net"
	"net/http"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/utils"

	"fmt"
	"log"
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

// waitForState waits for server reach certaing state
func waitForState(state RunningState) bool {
	retries := 10
	idle := 500 * time.Millisecond
	for ; retries > 0 && State != state; retries-- {
		time.Sleep(idle)
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
		return fmt.Errorf("Server is not in proper stopped state before trying to start it")
	}

	if utils.RunningOn(addr) {
		return fmt.Errorf("Another instance is running in same address %v", addr)
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
		var i int
		for i = 0; !utils.RunningOn(addr) && i < 10; i++ {
			time.Sleep(100 * time.Millisecond)
		}

		if i < 0 {
			log.Print("Server could not be started")
			return
		}

		State = Running
	}()

	// launching server in a goroutine for not blocking
	go func() {
		if listener != nil {
			err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
			if err != nil {
				log.Printf("HTTP Server closing - %v", err)
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
		return fmt.Errorf("Already stopped")
	}

	if listener == nil {
		State = Stopped
		return fmt.Errorf("Already closed")
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
