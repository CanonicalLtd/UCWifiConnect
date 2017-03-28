package netman

import (
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus"
)

type Client struct {
	dbusClient DbusClient
}

type DbusClient struct {
	test       bool
	BusObj     dbus.BusObject
	Connection *dbus.Conn
}

type Objecter interface {
	Object(dest string, path dbus.ObjectPath) dbus.BusObject
}

func (d *DbusClient) Object(dest string, path dbus.ObjectPath) dbus.BusObject {
	obj := d.Connection.Object(dest, path)
	return obj
}

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

func (c *Client) GetDevices() []string {
	c.dbusClient.Object("org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	setObject(c, "org.freedesktop.NetworkManager", "/org/freedesktop/NetworkManager")
	var devices []string
	err := c.dbusClient.BusObj.Call("org.freedesktop.NetworkManager.GetAllDevices", 0).Store(&devices)
	if err != nil {
		panic(err)
	}
	return devices
}

func (c *Client) GetWifiDevices(devices []string) []string {
	var wifiDevices []string
	for _, d := range devices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		deviceType, err2 := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.DeviceType")
		if err2 != nil {
			panic(err2)
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

func (c *Client) GetAccessPoints(devices []string, ap2device map[string]string) []string {
	var APs []string
	for _, d := range devices {
		var aps []string
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		err := c.dbusClient.BusObj.Call("org.freedesktop.NetworkManager.Device.Wireless.GetAllAccessPoints", 0).Store(&aps)
		if err != nil {
			panic(err)
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

type SSID struct {
	Ssid   string
	ApPath string
}

func (c *Client) GetSsids(APs []string, ssid2ap map[string]string) []SSID {
	var SSIDs []SSID
	for _, ap := range APs {
		objPath := dbus.ObjectPath(ap)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		ssid, err := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.AccessPoint.Ssid")
		if err != nil {
			panic(err)
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

func (c *Client) ConnectAp(ssid string, p string, ap2device map[string]string, ssid2ap map[string]string) {
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
}

func getSystemBus() *dbus.Conn {
	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to system bus:", err)
		panic(1)
	}
	return conn
}

func (c *Client) Ssids() ([]SSID, map[string]string, map[string]string) {
	ap2device := make(map[string]string)
	ssid2ap := make(map[string]string)
	devices := c.GetDevices()
	wifiDevices := c.GetWifiDevices(devices)
	APs := c.GetAccessPoints(wifiDevices, ap2device)
	SSIDs := c.GetSsids(APs, ssid2ap)
	return SSIDs, ap2device, ssid2ap
}

// check if no ethernet or wifi devices are connected
func (c *Client) Connected(devices []string) bool {
	for _, d := range devices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		dType, _ := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.DeviceType")
		state, _ := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.State")
		if dbus.Variant.Value(dType) != uint32(1) && dbus.Variant.Value(dType) != uint32(2) {
			continue
		}
		if dbus.Variant.Value(state) == uint32(100) {
			return true
		}
	}
	return false
}

// Check if connected by wifi
func (c *Client) ConnectedWifi(wifiDevices []string) bool {
	for _, d := range wifiDevices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		state, _ := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.State")
		if dbus.Variant.Value(state) == uint32(100) {
			return true
		}
	}
	return false
}

// Disconnects every interface passed. return shows number of disconnect calls  made
func (c *Client) DisconnectWifi(wifiDevices []string) int {
	ran := 0
	for _, d := range wifiDevices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		c.dbusClient.BusObj.Call("org.freedesktop.NetworkManager.Device.Disconnect", 0)
		ran += 1
	}
	return ran
}

// Set passed device to be managed by network manager, return iface set, if any
func (c *Client) SetIfaceManaged(iface string, devices []string) string {
	ran := ""
	for _, d := range devices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		iface_, err2 := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.Interface")
		if err2 != nil {
			fmt.Printf("Error 1 in SetIfaceManaged(): %v\n", err2)
			return ""
		}
		if iface != iface_.Value().(string) {
			continue
		}
		managed, err := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.Managed")
		if err != nil {
			fmt.Printf("Error 2 in SetIfaceManaged(): %v\n", err)
			return ""
		}
		if managed.Value().(bool) == true {
			return "" //no need to set, already managed
		}
		c.dbusClient.BusObj.Call("org.freedesktop.DBus.Properties.Set", 0, "org.freedesktop.NetworkManager.Device", "Managed", dbus.MakeVariant(true))
		ran = iface
		break
	}
	return ran
}

// Return  map[iface]device of wifi iterfaces that are managed by network manager
func (c *Client) WifisManaged(wifiDevices []string) (map[string]string, error) {
	ifaces := make(map[string]string)
	for _, d := range wifiDevices {
		objPath := dbus.ObjectPath(d)
		c.dbusClient.Object("org.freedesktop.NetworkManager", objPath)
		setObject(c, "org.freedesktop.NetworkManager", objPath)
		managed, err := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.Managed")
		if err != nil {
			fmt.Printf("Error 1 in WifisManaged(): %v\n", err)
			return ifaces, err
		}
		iface, err2 := c.dbusClient.BusObj.GetProperty("org.freedesktop.NetworkManager.Device.Interface")
		if err2 != nil {
			fmt.Printf("Error 2 in WifisManaged(): %v\n", err)
			return ifaces, err2
		}
		if managed.Value().(bool) == true {
			ifaces[iface.Value().(string)] = d
		}
	}
	return ifaces, nil
}
