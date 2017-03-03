package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
			fmt.Printf("====== device %s ap %s\n", d, i)
			APs = append(APs, i )
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

		lastSeen, err2 := obj.GetProperty("org.freedesktop.NetworkManager.AccessPoint.LastSeen")
		if err2 != nil {
			panic(err2)
		}
		fmt.Printf("====== %q LastSeen: %v\n", ssidStr, lastSeen)
		Ssid := SSID{ssid: ssidStr, apPath: ap}
		SSIDs = append(SSIDs, Ssid)
		//TODO: exclude ssid of device's own AP (the wifi-ap one)
	}
	return SSIDs
}

func connectAp(conn *dbus.Conn, ssid string, apPath string, devicePath string) {
/* a python example
connection_params = {
        "802-11-wireless": {
            "security": "802-11-wireless-security",
        },
        "802-11-wireless-security": {
            "key-mgmt": "wpa-psk",
            "psk": SEEKED_PASSPHRASE
        },
    }

    # Establish the connection.
    settings_path, connection_path = manager.AddAndActivateConnection(
        connection_params, device_path, our_ap_path)
    print "settings_path =", settings_path
    print "connection_path =", connection_path

heres the method sig:
AddAndActivateConnection (IN  a{sa{sv}} connection,
                          IN  o         device,
                          IN  o         specific_object,
                          OUT o         path,
                          OUT o         active_connection);

where (i think) a{sa{sv}} is
an array string:array maps of string: value maps
	var vals map[string]val
	var m[string]vals
	var []m

*/


	inner1 := make(map[string]dbus.Variant)
	inner1["security"] = dbus.MakeVariant("802-11-wireless-security")

	inner2 := make(map[string]dbus.Variant)
	inner2["key-mgmt"] = dbus.MakeVariant("wpa-psk")
	inner2["psk"] = dbus.MakeVariant("Redapple1")

	outer := make(map[string]map[string]dbus.Variant)
	outer["802-11-wireless"] = inner1
	outer["802-11-wireless-security"] = inner2
	fmt.Printf("%v\n",outer)
	obj := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	
	err := obj.Call("org.freedesktop.NetworkManager.AddAndActivateConnection", 0, outer, dbus.ObjectPath(devicePath), dbus.ObjectPath(apPath))
	if err != nil {
		fmt.Printf("===== activate call error: %v\n", err)
		panic(err)
	}
}

func main() {
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		os.Exit(1)
	}
	devices := getDevices(conn)
	wifiDevices := getWifiDevices(conn, devices)
	fmt.Printf("==== wifdevices: %v\n", wifiDevices)
	APs := getAccessPoints(conn, wifiDevices)
	fmt.Printf("==== APs: %v\n", APs)
	SSIDs := getSSIDs(conn, APs)
	for _, ssid := range SSIDs {
		fmt.Printf("%v\n", ssid)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("SSID: ")
	ssid, _ := reader.ReadString('\n')
	ssid = strings.TrimSpace(ssid)
	fmt.Print("Device path: ")
	dev, _ := reader.ReadString('\n')
	dev = strings.TrimSpace(dev)
	fmt.Print("AP path: ")
	ap, _ := reader.ReadString('\n')
	ap = strings.TrimSpace(ap)
	//connectAp(conn, "Mywifi", "/org/freedesktop/NetworkManager/AccessPoint/180", "/org/freedesktop/NetworkManager/Devices/2")
	connectAp(conn, "Mywifi", ap, dev)

	return

}
