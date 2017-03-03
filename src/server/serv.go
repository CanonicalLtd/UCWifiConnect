package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	//"github.com/godbus/dbus"
)

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}


/*
func getWifiNM() []string {

	bus = dbus.SystemBus()
	dbusObj = "/org/freedesktop/NetworkManager/Settings"
	//ussdstring = sys.argv[2]

	dbusIface = dbus.Interface(bus.get_object('org.freedesktop.NetworkManager', dbusObj),
                     'org.freedesktop.NetworkManager.Settings.ListConnections')
	fmt.Printf("connections:\n%v", dbusIface)
	var res []string
	//won't work need nmcli in the snap or user dbus
	//cmd := exec.Command("/snap/bin/network-manager.nmcli dev wifi list")
	return res
}
*/
func getWifi() []string {
	ldLibPath := os.Getenv("SNAP") + "/lib/arm-linux-gnueabihf"
	fmt.Printf("=== ldLibPath: %s\n", ldLibPath)
	os.Setenv("LD_LIBRARY_PATH", ldLibPath)
	wlanInterface := "wlan1"
	fmt.Println("==== 1")
	cmd := exec.Command(os.Getenv("SNAP") + "/sbin/iwlist",  wlanInterface, "scan")
	fmt.Println("==== 2")
	fmt.Println(cmd)
	fmt.Println("==== 3")
	//cmd := exec.Command("pwd")
	//cmd := exec.Command("/snap/bin/nmcli", "device", "wifi", "list")
	var out bytes.Buffer
	cmd.Stdout = &out
	fmt.Println("==== 4")
	err := cmd.Run()
	fmt.Println("==== 5")
	if err != nil {
		s := fmt.Sprintf("Run error. %s\n", err)
		return []string{s}
	}
	fmt.Println("==== 6")
	fmt.Printf("OUT: %q\n", out.String())
	fmt.Println("==== 7")
	lines := strings.Split(out.String(), "\n")
	var essids []string
	for _, l := range lines {
		if  strings.Contains(l, "ESSID") {
			parts := strings.Split(l, ":")
			if (len(parts[1]) < 3) { continue }
			e := parts[1][1:len(parts[1])-1]
			if contains(essids, e) { continue }
			essids = append(essids, e)
		}
	}
	return essids
	//return strings.Split(out.String(), "\n")

}

func para(s string) string {
	return fmt.Sprintf("<p>%s</p>", s)
}

func form(items []string) string {
	form_ := "<form>"
	for _, s := range items {
		line := fmt.Sprintf("<input type='radio' name='essidr' value='%s' checked>%s<br>", s, s)
		fmt.Println(line)
		form_ = form_ + line 
	}
	form_ = form_ + "</form>"
	return form_
}


func handler(w http.ResponseWriter, r *http.Request) {
	d, _ := os.Getwd()
	fmt.Fprintf(w, "<p>URL path: %s!</p>", r.URL.Path[1:])
	fmt.Fprintf(w, "<p>pwd: %s", d)
	//essids := getWifiNM()
	essids := getWifi()
	essids_form := form(essids)
	fmt.Fprintf(w, essids_form)
	//for _, s := range getWifi() {
	//	fmt.Fprintf(w, para(s))
	//}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
