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

// Network entry for a wifi network in page data
type Network struct {
	Ssid   string
	Status string
}

// PageData dynamic data to fulfill the template
type PageData struct {
	Networks []Network
	Alert    string
}

// SsidsHandler lists the current available SSIDs
func SsidsHandler(w http.ResponseWriter, r *http.Request) {

	var connectedWifi = ""
	if netman.ConnectedWifi() {
		//TODO get here the network we are connected to
		connectedWifi = "<TheWifiWeAreConnectedTo>"
	}

	// build dynamic data object
	ssids, _, _ := netman.Ssids()
	networks := make([]Network, len(ssids))
	for i := range ssids {
		status := "disconnected"
		if ssids[i].Ssid == connectedWifi {
			status = "connected"
		}
		networks[i] = Network{Ssid: ssids[i].Ssid, Status: status}
	}

	data := PageData{Networks: networks}

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

// ConnectHandler reads form got ssid and password and tries to connect to that network
func ConnectHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	ssids := r.Form["ssid"]
	if len(ssids) == 0 {
		log.Println("SSID not provided")
		return
	}
	ssid := ssids[0]

	pwd := ""
	pwds := r.Form["pwd"]
	if len(pwds) > 0 {
		pwd = pwds[0]
	}

	log.Printf("Connecting to %v...", ssid)

	//connect
	_, ap2device, ssid2ap := netman.Ssids()
	netman.ConnectAp(ssid, pwd, ap2device, ssid2ap)

	// redirect to list
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
