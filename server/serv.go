package main

import (
	"fmt"
	"net/http"

	"github.com/CanonicalLtd/UCWifiConnect/netman"
)

func para(s string) string {
	return fmt.Sprintf("<p>%s</p>", s)
}

func form(SSIDs []netman.SSID) string {
	form := "<form>"
	for _, s := range SSIDs {
		line := fmt.Sprintf("<input type='radio' name='essid' value='%s' checked>%s<br>", s.Ssid, s.Ssid)
		fmt.Println(line)
		form = form + line
	}
	form = form + "</form>"
	return form
}

func handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "<html><head></head></body>")
	ssids, _, _ := netman.Ssids()
	ssidsForm := form(ssids)
	fmt.Fprintf(w, ssidsForm)
	fmt.Fprintf(w, "</html>")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
