package server

import (
	"io"
	"net"
	"net/http"
	"time"

	"log"
)

const (
	address = ":8080"
)

// ListenAndServe starts http server.

// StartManagementServer starts web server for AP in management mode
/* Howto start and stop server
srvCLoser, err := server.StartManagementServer()
if err != nil {
	log.Fatalln("StartManagementServer Error - ", err)
}

// Do Stuff

// Close HTTP Server
err = srvCLoser.Close()
if err != nil {
	log.Fatalln("Server Close Error - ", err)
}

log.Println("Server Closed")
*/
func StartManagementServer() (sc io.Closer, err error) {
	return listenAndServe(address, managementHandler())
}

// StartExternalServer starts web server when connected to external WIFI
/* Howto start and stop server
srvCLoser, err := server.StartExternalServer()
if err != nil {
	log.Fatalln("StartExternalServer Error - ", err)
}

// Do Stuff

// Close HTTP Server
err = srvCLoser.Close()
if err != nil {
	log.Fatalln("Server Close Error - ", err)
}

log.Println("Server Closed")
*/
func StartExternalServer() (sc io.Closer, err error) {
	return listenAndServe(address, externalHandler())
}

type tcpKeepAliveListener struct {
	*net.TCPListener
}

// Accept accepts incomming tcp connections
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
