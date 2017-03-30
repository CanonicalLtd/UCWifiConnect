package main

import (
	"net/http"

	"log"

	"github.com/CanonicalLtd/UCWifiConnect/server"
	"github.com/gorilla/mux"
)

const (
	address = ":8080"
)

func handler() *mux.Router {
	router := mux.NewRouter()

	// Pages routes
	router.HandleFunc("/", server.SsidsHandler).Methods("GET")
	router.HandleFunc("/connect", server.ConnectHandler).Methods("POST")

	// Resources path
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(server.ResourcesPath)))
	router.PathPrefix("/static/").Handler(fs)

	return router
}

func main() {
	log.Fatal(http.ListenAndServe(address, handler()))
}
