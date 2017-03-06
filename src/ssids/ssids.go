package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus"
)

var ap2device map[string]string
var ssid2ap map[string]string

type options struct {
	getSsids bool
}

func getDevices(conn *dbus.Conn) []string {
	obj := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	var devices []string
	err2 := obj.Call("org.freedesktop.NetworkManager.GetAllDevices", 0).Store(&devices)
	if err2 != nil {
		panic(err2)
	}
	return devices
}

func getWifiDevices(conn *dbus.Conn, devices []string) []string {
	var wifiDevices [] string
	for _, d := range devices {
		objPath := dbus.ObjectPath(d)
		device := conn.Object("org.freedesktop.NetworkManager", objPath)
		deviceType, err2 := device.GetProperty("org.freedesktop.NetworkManager.Device.DeviceType")
		if err2 != nil {
			panic(err2)
		}
		var wifiType uint32
		wifiType = 2
		if deviceType.Value() != wifiType {
			continue
		}
		wifiDevices = append(wifiDevices, d)
	}
	return wifiDevices
}

func getAccessPoints(conn *dbus.Conn, devices []string) [] string {
	var APs [] string
	for _, d := range devices {
		objPath := dbus.ObjectPath(d)
		obj := conn.Object("org.freedesktop.NetworkManager", objPath)
		var aps []string
		err := obj.Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0).Store(&aps)
		if err != nil {
			panic(err)
		}
		for _, i := range aps {
			//fmt.Printf("====== device %s ap %s\n", d, i)
			APs = append(APs, i )
			ap2device[i] = d
		}
	}
	return APs
}

type SSID struct {
	ssid string
	apPath string
}

func getSSIDs(conn *dbus.Conn, APs []string) []SSID {
	var SSIDs []SSID
	//var SSIDs [] string
	for _, ap := range APs{
		objPath := dbus.ObjectPath(ap)
		obj := conn.Object("org.freedesktop.NetworkManager", objPath)
		ssid, err := obj.GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid")
		if err != nil {
			panic(err)
		}
		type  B []byte
		res := B(ssid.Value().([]byte))
		ssidStr := string(res)
		if len(ssidStr) < 1 {
			continue
		}
		found := false
		for _, s := range SSIDs {
			if s.ssid == ssidStr {
				found = true
			}
		}
		if found == true {
			continue
		}

		Ssid := SSID{ssid: ssidStr, apPath: ap}
		SSIDs = append(SSIDs, Ssid)
		ssid2ap[strings.TrimSpace(ssidStr)] = ap
		//TODO: exclude ssid of device's own AP (the wifi-ap one)
	}
	return SSIDs
}

func connectAp(conn *dbus.Conn, ssid string, p string) {
	//fmt.Println("OUT ssid|" + ssid + "|ssid")
	//fmt.Printf("OUR ssid to OUR ap:\n %v\n", ssid2ap[ssid])
	//fmt.Printf("ap to device:\n%v\n", ap2device)
	//fmt.Printf("ssid to ap:\n %v\n", ssid2ap)

	inner1 := make(map[string]dbus.Variant)
	inner1["security"] = dbus.MakeVariant("802-11-wireless-security")

	inner2 := make(map[string]dbus.Variant)
	inner2["key-mgmt"] = dbus.MakeVariant("wpa-psk")
	inner2["psk"] = dbus.MakeVariant(p)

	outer := make(map[string]map[string]dbus.Variant)
	outer["802-11-wireless"] = inner1
	outer["802-11-wireless-security"] = inner2
	fmt.Printf("%v\n",outer)

	fmt.Printf("dev path: %s\n",ap2device[ssid2ap[ssid]])
	fmt.Printf("ap path: %s\n",ssid2ap[ssid])

	obj := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	resp := obj.Call("org.freedesktop.NetworkManager.AddAndActivateConnection", 0, outer, dbus.ObjectPath(ap2device[ssid2ap[ssid]]), dbus.ObjectPath(ssid2ap[ssid]))
	fmt.Printf("===== activate call response:\n%v\n", resp)
}
func args() *options {
	opts := &options{}
	flag.BoolVar(&opts.getSsids, "get-ssids", false, "Connect to an AP")
	flag.Parse()
	return opts
}

func main() {
	opts := args()

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}

	ap2device = make(map[string]string)
	ssid2ap = make(map[string]string)

	devices := getDevices(conn)
	wifiDevices := getWifiDevices(conn, devices)
	//fmt.Printf("==== wifdevices: %v\n", wifiDevices)
	APs := getAccessPoints(conn, wifiDevices)
	//fmt.Printf("==== APs: %v\n", APs)
	SSIDs := getSSIDs(conn, APs)
	//fmt.Println("Found SSIDs:")
	if opts.getSsids {
		var out string
		for _, ssid := range SSIDs {
			out += strings.TrimSpace(ssid.ssid) + ","
		}
		fmt.Printf("%s\n", out[:len(out)-1])
		return
	}
	for _, ssid := range SSIDs {
		fmt.Printf("    %v\n", ssid.ssid)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Connect to SSID: ")
	ssid, _ := reader.ReadString('\n')
	ssid = strings.TrimSpace(ssid)
	fmt.Print("PW: ")
	pw, _ := reader.ReadString('\n')
	pw = strings.TrimSpace(pw)
	connectAp(conn, ssid, pw)

	return

}
