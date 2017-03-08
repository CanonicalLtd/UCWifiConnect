package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/knitzsche/wifi-connect/ssids"

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

func getWifi() []string {
	SSIDs, ap2device, ssid2ap := netman.Ssids()
	return SSIDs, ap2device, ssid2app
	/*
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
*/

}

func para(s string) string {
	return fmt.Sprintf("<p>%s</p>", s)
}

func form(SSIDs []string) string {
	form_ := "<form>"
	for _, s := range SSIDs {
		line := fmt.Sprintf("<input type='radio' name='essid' value='%s' checked>%s<br>", s.ssid, s.ssid)
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
	essids, _, _ := getWifi()
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
