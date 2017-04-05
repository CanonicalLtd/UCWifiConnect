package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// managementHandler handles requests for web UI when AP is up
func managementHandler() *mux.Router {
	router := mux.NewRouter()

	// Pages routes
	router.HandleFunc("/", SsidsHandler).Methods("GET")
	router.HandleFunc("/connect", ConnectHandler).Methods("POST")

	// Resources path
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(ResourcesPath)))
	router.PathPrefix("/static/").Handler(fs)

	return router
}

// operationalHandler handles request for web UI when connected to external WIFI
func operationalHandler() *mux.Router {

	//TODO IMPLEMENT
	return nil
}
