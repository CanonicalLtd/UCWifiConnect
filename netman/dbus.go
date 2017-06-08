/*
 * Copyright (C) 2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package netman

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus"
)

// Client type to support unit test mock and runtime execution
type Client struct {
	dbusClient DbusClient
}

// DbusClient properties for testing & runtime
type DbusClient struct {
	test       bool
	BusObj     dbus.BusObject
	Connection *dbus.Conn
}

// Objecter allows mocking the godbus Object function
type Objecter interface {
	Object(dest string, path dbus.ObjectPath) dbus.BusObject
}

// Object is the mock implementation of the godbus Object function
func (d *DbusClient) Object(dest string, path dbus.ObjectPath) dbus.BusObject {
	obj := d.Connection.Object(dest, path)
	return obj
}

// DefaultClient is the runtime client object
func DefaultClient() *Client {
	conn := getSystemBus()
	obj := conn.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	return &Client{
		dbusClient: DbusClient{
			test:       false,
			BusObj:     obj,
			Connection: getSystemBus(),
		},
	}
}

// NewClient is the mocked client object
func NewClient(myobj dbus.BusObject) *Client {
	return &Client{
		dbusClient: DbusClient{
			test:   true,
			BusObj: myobj,
		},
	}
}

// for Non test operation, save the current bus object to the client
func setObject(c *Client, iface string, path dbus.ObjectPath) {
	if !c.dbusClient.test {
		c.dbusClient.BusObj = getSystemBus().Object(iface, path)
	}
}

// GetDevices returns NetMan (NetworkManager) devices
func (c *Client) GetDevices() []string {
	if !c.dbusClient.test {
		c.dbusClient.Connection = getSystemBus()
	}
	c.dbusClient.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	setObject(c, "org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	var devices []string
	err := c.dbusClient.BusObj.Call("org.freedesktop.NetworkManager.GetAllDevices", 0).Store(&devices)
	if err != nil {
		fmt.Println("== wifi-connect: Error getting devices:", err)
	}
	return devices
}

// GetWifiDevices returns wifi NetMan devices
func (c *Client) GetWifiDevices(devices []string) []string {
	var wifiDevices []string
	for _, d := range devices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		deviceType, err2 := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.DeviceType")
		if err2 != nil {
			fmt.Println("== wifi-connect: Error getting wifi devices:", err2)
			continue
		}
		var wifiType uint32
		wifiType = 2
		if deviceType.Value() == nil {
			break
		}
		if deviceType.Value() != wifiType {
			continue
		}
		wifiDevices = append(wifiDevices, d)
	}
	return wifiDevices
}

//GetAccessPoints returns NetMan known external APs
func (c *Client) GetAccessPoints(devices []string, ap2device map[string]string) []string {
	var APs []string
	for _, d := range devices {
		var aps []string
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		err := c.dbusClient.BusObj.Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0).Store(&aps)
		if err != nil {
			fmt.Println("== wifi-connect: Error getting accesspoints:", err)
			continue
		}
		if len(aps) == 0 {
			break
		}
		for _, i := range aps {
			APs = append(APs, i)
			ap2device[i] = d
		}
	}
	return APs
}

// SSID holds SSID properties
type SSID struct {
	Ssid   string
	ApPath string
}

// getSsids returns known NetMan SSIDs
func (c *Client) getSsids(APs []string, ssid2ap map[string]string) []SSID {
	var SSIDs []SSID
	for _, ap := range APs {
		objPath := dbus.ObjectPath(ap)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		ssid, err := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid")
		if err != nil {
			fmt.Println("== wifi-connect: Error getting accesspoint's ssids:", err)
			continue
		}
		type B []byte
		res := B(ssid.Value().([]byte))
		ssidStr := string(res)
		if len(ssidStr) < 1 {
			continue
		}
		for _, s := range SSIDs {
			if s.Ssid == ssidStr {
				continue
			}
		}
		Ssid := SSID{Ssid: ssidStr, ApPath: ap}
		SSIDs = append(SSIDs, Ssid)
		ssid2ap[strings.TrimSpace(ssidStr)] = ap
		//TODO: exclude ssid of device's own AP (the wifi-ap one)
	}
	return SSIDs
}

// ConnectAp attempts to Connect to an external AP
func (c *Client) ConnectAp(ssid string, p string, ap2device map[string]string, ssid2ap map[string]string) error {
	inner1 := make(map[string]dbus.Variant)
	inner1["security"] = dbus.MakeVariant("802-11-wireless-security")

	inner2 := make(map[string]dbus.Variant)
	inner2["key-mgmt"] = dbus.MakeVariant("wpa-psk")
	inner2["psk"] = dbus.MakeVariant(p)

	outer := make(map[string]map[string]dbus.Variant)
	outer["802-11-wireless"] = inner1
	outer["802-11-wireless-security"] = inner2

	c.dbusClient.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	setObject(c, "org.freedesktop.NetworkManager", dbus.ObjectPath("/org/freedesktop/NetworkManager"))
	c.dbusClient.BusObj.Call("org.freedesktop.NetworkManager.AddAndActivateConnection", 0, outer, dbus.ObjectPath(ap2device[ssid2ap[ssid]]), dbus.ObjectPath(ssid2ap[ssid]))

	// loop until connected or until max loops
	trying := true
	idx := -1
	for trying {
		idx++
		time.Sleep(1000 * time.Millisecond)
		if c.Connected(c.GetWifiDevices(c.GetDevices())) {
			return nil
		}
		if idx == 19 {
			return errors.New("wifi-connect: cannot connect to AP")
		}
	}
	return nil
}

func getSystemBus() *dbus.Conn {
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "== wifi-connect: Error: Failed to connect to system bus:", err)
		panic(1)
	}
	return conn
}

// Ssids returns known SSIDs
func (c *Client) Ssids() ([]SSID, map[string]string, map[string]string) {
	ap2device := make(map[string]string)
	ssid2ap := make(map[string]string)
	devices := c.GetDevices()
	wifiDevices := c.GetWifiDevices(devices)
	APs := c.GetAccessPoints(wifiDevices, ap2device)
	SSIDs := c.getSsids(APs, ssid2ap)
	return SSIDs, ap2device, ssid2ap
}

// Connected checks if any passed ethernet/wifi devices are connected
func (c *Client) Connected(devices []string) bool {
	for _, d := range devices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		dType, err := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.DeviceType")
		if err != nil {
			fmt.Println("== wifi-connect: Error getting device type:", err)
			continue
		}
		state, err2 := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.State")
		if err2 != nil {
			fmt.Println("== wifi-connect: Error getting device state:", err2)
			continue
		}
		// only handle eth and wifi device type
		if dbus.Variant.Value(dType) != uint32(1) && dbus.Variant.Value(dType) != uint32(2) {
			continue
		}
		if dbus.Variant.Value(state) == uint32(100) {
			return true
		}
	}
	return false
}

// ConnectedWifi checks if any passed wifi devices are connected
func (c *Client) ConnectedWifi(wifiDevices []string) bool {
	for _, d := range wifiDevices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		state, err := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.State")
		if err != nil {
			fmt.Println("== wifi-connect: Error getting device state:", err)
			continue
		}
		if dbus.Variant.Value(state) == uint32(100) {
			return true
		}
	}
	return false
}

// DisconnectWifi disconnects every interface passed. return shows number of disconnect calls  made
func (c *Client) DisconnectWifi(wifiDevices []string) int {
	ran := 0
	for _, d := range wifiDevices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		c.dbusClient.BusObj.Call("org.freedesktop.NetworkManager.Device.Disconnect", 0)
		ran++
	}
	return ran
}

// SetIfaceManaged sets passed device to be managed/unmanaged by network manager, return iface set, if any
func (c *Client) SetIfaceManaged(iface string, state bool, devices []string) string {
	for _, d := range devices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		intface, err2 := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.Interface")
		if err2 != nil {
			fmt.Printf("== wifi-connect: Error in SetIfaceManaged() geting interface: %v\n", err2)
			return ""
		}
		if iface != intface.Value().(string) {
			continue
		}
		managed, err := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.Managed")
		if err != nil {
			fmt.Printf("== wifi-connect: Error in SetIfaceManaged() fetching device managed: %v\n", err)
			return ""
		}
		switch state {
		case true:
			if managed.Value().(bool) == true {
				return "" //no need to set, already managed
			}
		case false:
			if managed.Value().(bool) == false {
				return "" //no need to set, already UNmanaged
			}
		}

		c.dbusClient.BusObj.Call("org.freedesktop.DBus.Properties.Set", 0, "org.freedesktop.NetworkManager.Device", "Managed", dbus.MakeVariant(state))
		// loop until interface is in desired managed state or max iters reached
		idx := -1
		for {
			idx++
			time.Sleep(1000 * time.Millisecond)
			managedState, err := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.State")
			if err == nil {
				switch state {
				case true:
					if managedState.Value() == uint32(30) { //NM_DEVICE_STATE_DISCONNECTED
						return iface
					}
				case false:
					if managedState.Value() == uint32(10) { //NM_DEVICE_STATE_UNMANAGED
						return iface
					}
				}
			}
			if idx == 59 { //give it 60 iters ~= one minute
				break
			}
		}
	}
	return "" //no iface state changed
}

// WifisManaged returns  map[iface]device of wifi iterfaces that are managed by network manager
func (c *Client) WifisManaged(wifiDevices []string) (map[string]string, error) {
	ifaces := make(map[string]string)
	for _, d := range wifiDevices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		managed, err := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.Managed")
		if err != nil {
			fmt.Printf("== wifi-connect: Error in WifisManaged() getting device managed : %v\n", err)
			return ifaces, err
		}
		iface, err2 := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.Interface")
		if err2 != nil {
			fmt.Printf("== wifi-connect: Error in WifisManaged() getting device interface: %v\n", err)
			return ifaces, err2
		}
		if managed.Value().(bool) == true {
			ifaces[iface.Value().(string)] = d
		}
	}
	return ifaces, nil
}
