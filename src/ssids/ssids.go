package main

import (
	"fmt"
	"os"

	"github.com/godbus/dbus"
)

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
			APs = append(APs, i )
		}
	}
	return APs
}
func getSSIDs(conn *dbus.Conn, APs []string) [] string {
	var SSIDs [] string
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
			if s == ssidStr {
				found = true
			}
		}
		if found == false {
			SSIDs = append(SSIDs, ssidStr)
		}
		//TODO: exclude ssid of device's own AP (the wifi-ap one)
	}
	return SSIDs
}
func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	devices := getDevices(conn)
	wifiDevices := getWifiDevices(conn, devices)
	//TODO RequestScan() to ensure up to date, unless this ttakes too long
	//fmt.Printf("==== wifdevices: %v\n", wifiDevices)
	APs := getAccessPoints(conn, wifiDevices)
	//fmt.Printf("==== APs: %v\n", APs)
	SSIDs := getSSIDs(conn, APs)
	for _, ssid := range SSIDs { 	
		fmt.Println(ssid)
	}

}
