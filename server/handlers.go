package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
)

const (
	templatePath = "/templates/ssids.html"
)

// ResourcesPath absolute path to web static resources
var ResourcesPath = filepath.Join(os.Getenv("SNAP"), "static")

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
	templateAbsPath := filepath.Join(ResourcesPath, templatePath)
	t, err := template.ParseFiles(templateAbsPath)
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
