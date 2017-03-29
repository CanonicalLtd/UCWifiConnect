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
	Networks []netman.SSID
}

// SsidsHandler lists the current available SSIDs
func SsidsHandler(w http.ResponseWriter, r *http.Request) {
	c := netman.DefaultClient()
	// build dynamic data object
	networks, _, _ := c.Ssids()
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
	c := netman.DefaultClient()
	_, ap2device, ssid2ap := c.Ssids()
	c.ConnectAp(ssid, pwd, ap2device, ssid2ap)

	// redirect to result web
	http.Redirect(w, r, "/static/templates/connect_result.html", http.StatusMovedPermanently)
}
