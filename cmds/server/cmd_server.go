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
	router.HandleFunc("/", server.SsidsHandler).Methods("GET")
	return router
}

func main() {
	log.Fatal(http.ListenAndServe(address, handler()))
}
