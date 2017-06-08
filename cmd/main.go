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

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
	"github.com/CanonicalLtd/UCWifiConnect/server"
	"github.com/CanonicalLtd/UCWifiConnect/utils"
	"github.com/CanonicalLtd/UCWifiConnect/wifiap"

	"github.com/gorilla/mux"
)

func help() string {

	text :=
		`Usage: sudo wifi-connect COMMAND [VALUE]

Commands:
	stop:	 		Disables wifi-connect from automatic control, leaving system 
				in current state
	start:	 		Enables wifi-connect as automatic controller, restarting from
				a clean state
	show-ap:		Show AP configuration
	ssid VALUE: 		Set the AP ssid (causes AP restart if it is UP)
	passphrase VALUE: 	Set the AP passphrase (cause AP restart if it is UP)
`
	return text
}

func mgmtHandler() *mux.Router {
	router := mux.NewRouter()

	// Pages routes
	router.HandleFunc("/", server.ManagementHandler).Methods("GET")
	router.HandleFunc("/connect", server.ConnectHandler).Methods("POST")

	// Resources path
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(server.ResourcesPath)))
	router.PathPrefix("/static/").Handler(fs)

	return router
}

func operHandler() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/", server.OperationalHandler).Methods("GET")
	router.HandleFunc("/hashit", server.HashItHandler).Methods("POST")

	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(server.ResourcesPath)))
	router.PathPrefix("/static/").Handler(fs)

	return router
}

// checkSudo return false if the current user is not root, else true
func checkSudo() bool {
	if os.Geteuid() != 0 {
		fmt.Println("Error: This command requires sudo")
		return false
	}
	return true
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("== wifi-connect/cmd Error: no command arguments provided")
		return
	}
	args := os.Args[1:]

	switch args[0] {
	case "help":
		fmt.Printf("%s\n", help())
	case "-help":
		fmt.Printf("%s\n", help())
	case "-h":
		fmt.Printf("%s\n", help())
	case "--help":
		fmt.Printf("%s\n", help())
	case "stop":
		if !checkSudo() {
			return
		}
		err := utils.WriteFlagFile(os.Getenv("SNAP_COMMON") + "/manualMode")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Entering MANUAL Mode. Wifi-connect has stopped managing state. Use 'start' to restore normal operations")
	case "start":
		if !checkSudo() {
			return
		}
		err := utils.RemoveFlagFile(os.Getenv("SNAP_COMMON") + "/manualMode")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Entering NORMAL Mode.")
	case "show-ap":
		if !checkSudo() {
			return
		}
		wifiAPClient := wifiap.DefaultClient()
		result, err := wifiAPClient.Show()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if result != nil {
			utils.PrintMapSorted(result)
			return
		}
	case "ssid":
		if !checkSudo() {
			return
		}
		if len(os.Args) < 3 {
			fmt.Println("Error: no ssid provided")
			return
		}
		wifiAPClient := wifiap.DefaultClient()
		wifiAPClient.SetSsid(os.Args[2])
	case "passphrase":
		if !checkSudo() {
			return
		}
		if len(os.Args) < 3 {
			fmt.Println("Error: no passphrase provided")
			return
		}
		if len(os.Args[2]) < 13 {
			fmt.Println("Error: passphrase must be at least 13 chars long")
			return
		}
		wifiAPClient := wifiap.DefaultClient()
		wifiAPClient.SetPassphrase(os.Args[2])
	case "get-devices":
		c := netman.DefaultClient()
		devices := c.GetDevices()
		for d := range devices {
			fmt.Println(d)
		}
	case "get-wifi-devices":
		c := netman.DefaultClient()
		devices := c.GetWifiDevices(c.GetDevices())
		for d := range devices {
			fmt.Println(d)
		}
	case "get-ssids":
		c := netman.DefaultClient()
		SSIDs, _, _ := c.Ssids()
		var out string
		for _, ssid := range SSIDs {
			out += strings.TrimSpace(ssid.Ssid) + ","
		}
		if len(out) > 0 {
			fmt.Printf("%s\n", out[:len(out)-1])
		}
	case "check-connected":
		c := netman.DefaultClient()
		if c.ConnectedWifi(c.GetWifiDevices(c.GetDevices())) {
			fmt.Println("Device is connected")
		} else {
			fmt.Println("Device is not connected")
		}

	case "check-connected-wifi":
		c := netman.DefaultClient()
		if c.ConnectedWifi(c.GetWifiDevices(c.GetDevices())) {
			fmt.Println("Device is connected to external wifi AP")
		} else {
			fmt.Println("Device is not connected to external wifi AP")
		}
	case "disconnect-wifi":
		c := netman.DefaultClient()
		c.DisconnectWifi(c.GetWifiDevices(c.GetDevices()))
	case "wifis-managed":
		c := netman.DefaultClient()
		wifis, err := c.WifisManaged(c.GetWifiDevices(c.GetDevices()))
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range wifis {
			fmt.Printf("%s : %s\n", k, v)
		}
	case "manage-iface":
		if len(os.Args) < 3 {
			fmt.Println("Error: no interface provided")
			return
		}
		c := netman.DefaultClient()
		c.SetIfaceManaged(os.Args[2], true, c.GetWifiDevices(c.GetDevices()))
	case "unmanage-iface":
		if len(os.Args) < 3 {
			fmt.Println("Error: no interface provided")
			return
		}
		c := netman.DefaultClient()
		c.SetIfaceManaged(os.Args[2], false, c.GetWifiDevices(c.GetDevices()))
	case "connect":
		c := netman.DefaultClient()
		SSIDs, ap2device, ssid2ap := c.Ssids()
		for _, ssid := range SSIDs {
			fmt.Printf("    %v\n", ssid.Ssid)
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Connect to AP. Enter SSID: ")
		ssid, _ := reader.ReadString('\n')
		ssid = strings.TrimSpace(ssid)
		fmt.Print("Enter phasprase: ")
		pw, _ := reader.ReadString('\n')
		pw = strings.TrimSpace(pw)
		c.ConnectAp(ssid, pw, ap2device, ssid2ap)
	case "management":
		http.ListenAndServe(":8081", mgmtHandler())
	case "operational":
		http.ListenAndServe(":8081", operHandler())
	case "set-portal-password":
		if len(os.Args) < 3 {
			fmt.Println("Error: no string to hash provided")
			return
		}
		if len(os.Args[2]) < 8 {
			fmt.Println("Error: password must be at least 8 characters long")
		}
		b, err := utils.HashIt(os.Args[2])
		if err != nil {
			fmt.Println("Error hashing:", err)
		}
		fmt.Println(string(b))
	default:
		fmt.Println("Error. Your command is not supported. Please try 'help'")
	}
}
