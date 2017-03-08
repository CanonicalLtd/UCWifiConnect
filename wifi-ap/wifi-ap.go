package wifi-ap

import (
	"fmt"
	"os"
	"os/exec"
)


func enable() {
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"disabled":"false"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err)
		return
	}
	return
}

func disable() {
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"disabled":"true"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err)
		return
	}
	return
}

func setSsid(ssid string) {
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"wifi.ssid":"` + ssid +`"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err)
		return
	}
	return
}
func setPassphrase(passphrase string) {
	//for now, let's use wpa2 security
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"wifi.security":"wpa2"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err)
		return
	}
	//passphrase
	path = os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err2 := exec.Command(path, "-d", `{"wifi.security-passphrase":"` + passphrase +`"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err2 != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err2)
		return
	}
	return
}
