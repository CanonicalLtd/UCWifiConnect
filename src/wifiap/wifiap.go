package wifiap

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Show shows current wifi-ap status
func Show() {
	result, err := DefaultRestClient().Show()
	if err != nil {
		log.Printf("wifi-ap show operation failed: %q\n", err)
		return
	}
	// TODO see if needed here to return value or simply showing that in stdout is ok
	printMapSorted(result)
}

// Enable enables wifi ap
func Enable() {
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"disabled":"false"}`, os.Getenv("SNAP_COMMON")+"/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path, err)
		return
	}
	return
}

// Disable disables wifi ap
func Disable() {
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"disabled":"true"}`, os.Getenv("SNAP_COMMON")+"/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path, err)
		return
	}
	return
}

// SetSsid sets the ssid for the wifi ap
func SetSsid(ssid string) {
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"wifi.ssid":"`+ssid+`"}`, os.Getenv("SNAP_COMMON")+"/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path, err)
		return
	}
	return
}

// SetPassphrase sets the credential to access the wifi ap
func SetPassphrase(passphrase string) {
	//for now, let's use wpa2 security
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"wifi.security":"wpa2"}`, os.Getenv("SNAP_COMMON")+"/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path, err)
		return
	}
	//passphrase
	path = os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err2 := exec.Command(path, "-d", `{"wifi.security-passphrase":"`+passphrase+`"}`, os.Getenv("SNAP_COMMON")+"/sockets/control", "/v1/configuration").Output()
	if err2 != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path, err2)
		return
	}
	return
}
