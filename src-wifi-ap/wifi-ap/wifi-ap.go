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
	verbose bool
	err           string
}
func args() *options {

	opts := &options{}
	flag.BoolVar(&opts.show, "show", false, "Show the wifi-ap confiruation")
	flag.BoolVar(&opts.enable, "ap-on", false, "Turn on the AP")
	flag.BoolVar(&opts.disable, "ap-off", false, "Turn off the AP")
	flag.BoolVar(&opts.verbose, "verbose", false, "Display verbose output")
	flag.Parse()
	return opts
}

func show() {
	cmdArgs := []string{os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration"}
	snapPath := os.Getenv("SNAP")
	path := snapPath + "/bin/unixhttpc"
	//fmt.Printf("SHOW cmd: %q \n", path)
	//fmt.Printf("SHOW args: %q\n", cmdArgs)
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
	_, err := exec.Command(path, "-d", `{"disabled":"itrue"}`, os.Getenv("SNAP_COMMON") + "/sockets/control", "/v1/configuration").Output()
	if err != nil {
		fmt.Printf("Error: '%q' failed. %q\n", path,  err)
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
	case opts.enable:
		enable(opts)
	case opts.disable:
		disable(opts)
	}
}
