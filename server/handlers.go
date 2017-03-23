package server

import (
	"log"
	"net/http"
	"text/template"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
)

const (
	templatePath = "templates/ssids.html"
)

// PageData dynamic data to fulfill the template
type PageData struct {
	Ssids []netman.SSID
}

// SsidsHandler lists the current available SSIDs
func SsidsHandler(w http.ResponseWriter, r *http.Request) {

	// build dynamic data object
	ssids, _, _ := netman.Ssids()
	data := PageData{ssids}

	// parse template
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Error loading the template at %v : %v\n", templatePath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("Error executing the template at %v : %v\n", templatePath, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
