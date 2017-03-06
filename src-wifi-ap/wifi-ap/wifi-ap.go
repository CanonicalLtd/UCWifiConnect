package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func contains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
type options struct {
	show bool
	enable bool
	disable bool
	ssid string
	passphrase string
	security string
	verbose bool
	err           string
}
func args() *options {

	opts := &options{}
	flag.BoolVar(&opts.show, "show", false, "Show the wifi-ap confiruation")
	flag.BoolVar(&opts.enable, "ap-on", false, "Turn on the AP")
	flag.BoolVar(&opts.disable, "ap-off", false, "Turn off the AP")
	flag.StringVar(&opts.ssid, "ssid", "", "Set the AP's SSID")
	flag.StringVar(&opts.passphrase, "passphrase", "", "Set the AP's passphrase")
	flag.BoolVar(&opts.verbose, "verbose", false, "Display verbose output")
	opts.security = "wpa2"
	flag.Parse()
	return opts
}

func show() {
	cmdArgs := []string{os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration"}
	snapPath := os.Getenv("SNAP")
	path := snapPath + "/bin/unixhttpc"
	out, err := exec.Command(path, cmdArgs...).Output()
	if err != nil {
		fmt.Printf("Error: '%q %q' failed. %q\n", path, cmdArgs, err)
		return
	}
	fmt.Printf("Wifi-ap Configuration:\n%s\n",out)
	return
}

func enable(opts * options) {
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"disabled":"false"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err)
		return
	}
	return
}

func disable(opts * options) {
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"disabled":"true"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err)
		return
	}
	return
}

func setSsid(opts * options) {
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"wifi.ssid":"` + opts.ssid +`"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err)
		return
	}
	return
}
func setPassphrase(opts * options) {
	//for now, let's use wpa2 security
	path := os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err := exec.Command(path, "-d", `{"wifi.security":"wpa2"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err)
		return
	}
	//passphrase
	path = os.Getenv("SNAP") + "/bin/unixhttpc"
	_, err2 := exec.Command(path, "-d", `{"wifi.security-passphrase":"` + opts.passphrase +`"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err2 != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err2)
		return
	}
	return
}
func main() {
	opts := args()
	if len(opts.err) > 0 {
		fmt.Printf("%q. Stopping.\n", opts.err)
		return
	}
	switch {
	case opts.show:
		show()
	case len(opts.ssid) > 1:
		setSsid(opts)
	case len(opts.passphrase) > 1:
		if len(opts.passphrase) < 13 {
			fmt.Println("Passphrase must be at least 13 chars in length. Please try again.")
			return
		}
		setPassphrase(opts)
	case opts.enable:
		enable(opts)
	case opts.disable:
		disable(opts)
	}
}
