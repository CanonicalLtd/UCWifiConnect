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
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
	"github.com/CanonicalLtd/UCWifiConnect/utils"
	"github.com/CanonicalLtd/UCWifiConnect/wifiap"
)

const (
	ssidsTemplatePath      = "/templates/ssids.html"
	connectingTemplatePath = "/templates/connecting.html"
)

// ResourcesPath absolute path to web static resources
var ResourcesPath = filepath.Join(os.Getenv("SNAP"), "static")

// SsidsFile path to the file filled by daemon with available ssids in csv format
var SsidsFile = filepath.Join(os.Getenv("SNAP_COMMON"), "ssids")

// Data interface representing any data included in a template
type Data interface{}

// SsidsData dynamic data to fulfill the SSIDs page template
type SsidsData struct {
	Ssids []string
}

// ConnectingData dynamic data to fulfill the connect result page template
type ConnectingData struct {
	Ssid string
}

func execTemplate(w http.ResponseWriter, templatePath string, data Data) {
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

func readSsidsFile() ([]string, error) {
	f, err := os.Open(SsidsFile)
	if err != nil {
		log.Printf("Error:%v", err)
		return nil, err
	}

	reader := csv.NewReader(f)
	// all ssids are in the same record
	record, err := reader.Read()
	if err == io.EOF {
		empty := make([]string, 0)
		return empty, nil
	}
	return record, err
}

// SsidsHandler lists the current available SSIDs
func SsidsHandler(w http.ResponseWriter, r *http.Request) {
	// daemon stores current available ssids in a file
	ssids, err := readSsidsFile()
	if err != nil {
		log.Printf("Error reading SSIDs file: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := SsidsData{Ssids: ssids}

	// parse template
	execTemplate(w, ssidsTemplatePath, data)
}

// ConnectHandler reads form got ssid and password and tries to connect to that network
func ConnectHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	ssids := r.Form["ssid"]
	if len(ssids) == 0 {
		log.Println("== wifi-connect: SSID not provided")
		return
	}
	ssid := ssids[0]

	data := ConnectingData{ssid}
	execTemplate(w, connectingTemplatePath, data)

	pwd := ""
	pwds := r.Form["pwd"]
	if len(pwds) > 0 {
		pwd = pwds[0]
	}

	log.Printf("== wifi-connect: Connecting to %v...", ssid)

	cw := wifiap.DefaultClient()
	cw.Disable()

	//give time for wifi-ap to complete shudown
	time.Sleep(50000 * time.Millisecond)

	//connect
	c := netman.DefaultClient()
	c.SetIfaceManaged("wlan0", true, c.GetWifiDevices(c.GetDevices()))
	_, ap2device, ssid2ap := c.Ssids()

	c.ConnectAp(ssid, pwd, ap2device, ssid2ap)

	//wait, to provide time for the connection to occur
	time.Sleep(30000 * time.Millisecond)

	//remove flag file so that daemon starts checking state
	//and takes control again
	utils.RemoveWaitFile()
}
