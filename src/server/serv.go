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
	
	var essids []string
	cmd := exec.Command(os.Getenv("SNAP") + "/bin/ssids",  "-get-ssids")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		s := fmt.Sprintf("Run error. %s\n", err)
		return []string{s}
	}
	fmt.Printf("SSIDs: %q\n", out.String())
	res := strings.TrimSpace(out.String())
	essids = strings.Split(res, ",")
	return essids

}

func para(s string) string {
	return fmt.Sprintf("<p>%s</p>", s)
}

func form(items []string) string {
	form_ := "<form>"
	for _, s := range items {
		line := fmt.Sprintf("<input type='radio' name='essid' value='%s' checked>%s<br>", s, s)
		fmt.Println(line)
		form_ = form_ + line 
	}
	form_ = form_ + "</form>"
	return form_
}

func handler(w http.ResponseWriter, r *http.Request) {
	d, _ := os.Getwd()
	//fmt.Fprintf(w, "<p>URL path: %s</p>", r.URL.Path[1:])
	fmt.Fprintf(w, "<p>pwd: %s", d)
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
