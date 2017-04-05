package server

import (
	"io"
	"net"
	"net/http"
	"time"

	"log"
)

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts incoming tcp connections
func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func listenAndServe(addr string, handler http.Handler) (sc io.Closer, err error) {

	var listener net.Listener

	srv := &http.Server{Addr: addr, Handler: handler}

	listener, err = net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	// launching server in a goroutine for not blocking
	go func() {
		err := srv.Serve(tcpKeepAliveListener{listener.(*net.TCPListener)})
		if err != nil {
			log.Println("HTTP Server Error - ", err)
		}
	}()

	return listener, nil
}
