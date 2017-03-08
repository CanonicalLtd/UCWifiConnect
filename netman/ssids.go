package ssids

import (
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
	var wifiDevices []string
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

func getAccessPoints(conn *dbus.Conn, devices []string, ap2device map[string]string) [] string {
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
			ap2device[i] = d
		}
	}
	return APs
}

type SSID struct {
	Ssid string
	ApPath string
}

func getSSIDs(conn *dbus.Conn, APs []string, ssid2ap map[string]string) []SSID {
	var SSIDs []SSID
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
			if s.Ssid == ssidStr {
				found = true
			}
		}
		if found == true {
			continue
		}

		Ssid := SSID{Ssid: ssidStr, ApPath: ap}
		SSIDs = append(SSIDs, Ssid)
		ssid2ap[strings.TrimSpace(ssidStr)] = ap
		//TODO: exclude ssid of device's own AP (the wifi-ap one)
	}
	return SSIDs
}

func ConnectAp(ssid string, p string, ap2device map[string]string, ssid2ap map[string]string) {
	conn := getSystemBus()
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
	obj.Call("org.freedesktop.NetworkManager.AddAndActivateConnection", 0, outer, dbus.ObjectPath(ap2device[ssid2ap[ssid]]), dbus.ObjectPath(ssid2ap[ssid]))
	//fmt.Printf("===== activate call response:\n%v\n", resp)
}
func getSystemBus() *dbus.Conn {
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		panic(1)
	}
	return conn
}

func Ssids() ([]SSID, map[string]string, map[string]string) {
	conn := getSystemBus()

	ap2device := make(map[string]string)
	ssid2ap := make(map[string]string)

	devices := getDevices(conn)
	wifiDevices := getWifiDevices(conn, devices)
	APs := getAccessPoints(conn, wifiDevices, ap2device)
	SSIDs := getSSIDs(conn, APs, ssid2ap)
	return SSIDs, ap2device, ssid2ap
}
