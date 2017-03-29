package netman

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/godbus/dbus"
)

// Make a dbus Call object for testing
func makeCall() *dbus.Call {
	objPath := dbus.ObjectPath("objpath")
	args2 := make([]interface{}, 3)
	devices := []string{"/d/1", "/d/2"}
	body := []interface{}{devices}
	call := &dbus.Call{
		Destination: "dest",
		Path:        objPath,
		Method:      "method",
		Args:        args2,
		//Done:        channel,
		Err:  nil,
		Body: body,
	}
	return call
}

type mockObj struct {
	wifiDevices []string
	ssids       [][]byte
	aps         []string
	ifaces      []string
	managed     bool
}

func (mock *mockObj) Object(dest string, path dbus.ObjectPath) dbus.BusObject {
	return mock
}

func (mock *mockObj) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {
	call := makeCall()
	switch method {
	case "org.freedesktop.NetworkManager.GetAllDevices":
		devices := []string{"/d/1", "/d/2", "/d/3"}
		body := []interface{}{devices}
		call.Body = body
	case "org.freedesktop.NetworkManager.GetAllAccessPoints":
		ap1 := "/ap/" + strconv.Itoa(len(mock.aps)+1)
		ap2 := "/ap/" + strconv.Itoa(len(mock.aps)+2)
		aps := []string{ap1, ap2}
		mock.aps = append(mock.aps, ap1)
		mock.aps = append(mock.aps, ap2)
		body := []interface{}{aps}
		call.Body = body
	case "org.freedesktop.NetworkManager.Device.Disconnect":
	}
	//fmt.Println("mockCall body: ", call.Body)
	return call
}

func (mock *mockObj) Go(method string, flags dbus.Flags, ch chan *dbus.Call, args ...interface{}) *dbus.Call {
	call := makeCall()
	return call
}

func (mock *mockObj) GetProperty(p string) (dbus.Variant, error) {
	switch p {
	case "org.freedesktop.NetworkManager.Device.DeviceType":
		if len(mock.wifiDevices) < 2 { // 2 of three devices are wifi
			mock.wifiDevices = append(mock.wifiDevices, "wifi"+strconv.Itoa(len(mock.wifiDevices)))
			return dbus.MakeVariant(uint32(2)), nil
		}
		return dbus.MakeVariant(uint32(1)), nil
	case "org.freedesktop.NetworkManager.AccessPoint.Ssid":
		ssid := "ssid" + strconv.Itoa(len(mock.ssids))
		ssidB := []byte(ssid)
		mock.ssids = append(mock.ssids, ssidB)
		return dbus.MakeVariant(ssidB), nil
	case "org.freedesktop.NetworkManager.Device.State":
		return dbus.MakeVariant(uint32(100)), nil
	case "org.freedesktop.NetworkManager.Device.Managed":
		if mock.managed {
			return dbus.MakeVariant(true), nil
		}
		return dbus.MakeVariant(false), nil
	case "org.freedesktop.NetworkManager.Device.Interface":
		if len(mock.ifaces) == 0 {
			mock.ifaces = append(mock.ifaces, "iface0")
			return dbus.MakeVariant("iface0"), nil
		}
		if len(mock.ifaces) == 1 {
			mock.ifaces = append(mock.ifaces, "iface1")
			return dbus.MakeVariant("iface1"), nil
		}
	}
	return dbus.MakeVariant("GetProperty error"), errors.New("no such property found")
}

func (mock *mockObj) Destination() string {
	return "destination"
}

func (mock *mockObj) Path() dbus.ObjectPath {
	return dbus.ObjectPath("/fake/objectPath")
}

func (mock *mockObj) systemBus() (*dbus.Conn, error) {
	conn := &dbus.Conn{}
	return conn, nil
}

func TestGetDevices(t *testing.T) {
	client := NewClient(&mockObj{})
	devices := client.GetDevices()
	found1 := false
	found2 := false
	found3 := false
	for _, v := range devices {
		switch v {
		case "/d/1":
			found1 = true
		case "/d/2":
			found2 = true
		case "/d/3":
			found3 = true
		}
	}
	if !found1 || !found2 || !found3 {
		t.Errorf("An expected device was not found")
	}
	fmt.Printf("===== Found devices: %v\n", devices)
}

func TestGetWifiDevices(t *testing.T) {
	client := NewClient(&mockObj{})
	devices := client.GetDevices()
	wifiDevices := client.GetWifiDevices(devices)
	found1 := false
	found2 := false
	for _, v := range devices {
		switch v {
		case "/d/1":
			found1 = true
		case "/d/2":
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Errorf("An expected device was not found")
	}
	if len(wifiDevices) != 2 {
		t.Errorf("Two wifi device should have been found but, found: %d", len(wifiDevices))
	}
	fmt.Printf("===== Found wifi devices: %v\n", wifiDevices)
}

func TestGetAPs(t *testing.T) {
	client := NewClient(&mockObj{})
	devices := client.GetDevices()
	ap2device := make(map[string]string)
	wifiDevices := client.GetWifiDevices(devices)
	aps := client.GetAccessPoints(wifiDevices, ap2device)
	if len(aps) != 4 {
		t.Errorf("4 APs  should have been found, but found: %d", len(aps))
	}
	fmt.Printf("===== Found APs: %v\n", aps)
}

func TestGetSsids(t *testing.T) {
	client := NewClient(&mockObj{})
	devices := client.GetDevices()
	wifiDevices := client.GetWifiDevices(devices)
	ap2device := make(map[string]string)
	ssid2ap := make(map[string]string)
	aps := client.GetAccessPoints(wifiDevices, ap2device)
	ssids := client.getSsids(aps, ssid2ap)
	if len(ssids) != 4 {
		t.Errorf("4 SSIDs should have been found, but found: %d", len(ssids))
	}
	fmt.Printf("===== GetSSIDs (ssid/ap): %v\n", ssids)
}

func TestSsids(t *testing.T) {
	client := NewClient(&mockObj{})
	ssids, _, _ := client.Ssids()
	if len(ssids) != 4 {
		t.Errorf("4 SSIDs should have been found, but found: %d", len(ssids))
	}
	fmt.Printf("===== Ssids() (ssid/ap): %v\n", ssids)
}

func TestConnected(t *testing.T) {
	client := NewClient(&mockObj{})
	if !client.Connected([]string{"d1"}) {
		t.Errorf("Should have found connected state, but did not")
	}
	if client.Connected([]string{}) {
		t.Errorf("Should have found no connection since there are no devices, but did not")
	}
}

func TestConnectedWifi(t *testing.T) {
	client := NewClient(&mockObj{})
	if !client.ConnectedWifi([]string{"d1"}) {
		t.Errorf("Should have found Wificonnected state, but did not")
	}
	if client.ConnectedWifi([]string{}) {
		t.Errorf("Should have found no connection since there are no devices, but did not")
	}
}

func TestiDiscconnectWifi(t *testing.T) {
	client := NewClient(&mockObj{})
	res := client.DisconnectWifi([]string{})
	if res != 0 {
		t.Errorf("0 Disconnect call expected, but found: %d", res)
	}
	res = client.DisconnectWifi([]string{"d1"})
	if res != 1 {
		t.Errorf("1 Disconnect call expected, but found: %d", res)
	}
}

func TestSetIfaceManaged(t *testing.T) {
	mock := &mockObj{}
	client := NewClient(mock)
	res := client.SetIfaceManaged("notaniface", []string{})
	if res != "" {
		t.Errorf("No Disconnect call expected, but found: %s", res)
	}
	res = client.SetIfaceManaged("iface2", []string{"d0", "d1"})
	if res != "" {
		t.Errorf("No Disconnect call expected, but found: %s", res)
	}
	mock.ifaces = mock.ifaces[:0]
	res = client.SetIfaceManaged("iface0", []string{"d0"})
	if res != "iface0" {
		t.Errorf("Disconnect call to iface0 expected, but found: %s", res)
	}
	mock.ifaces = mock.ifaces[:0]
	res = client.SetIfaceManaged("iface1", []string{"d0", "d1"})
	if res != "iface1" {
		t.Errorf("Disconnect call to iface1 expected, but found: %s", res)
	}
	mock.ifaces = mock.ifaces[:0]
	mock.managed = true
	res = client.SetIfaceManaged("iface1", []string{"d0", "d1", "d3"})
	if res != "" {
		t.Errorf("Disconnect call not expected (all ifaces managed), but found: %s", res)
	}
}

func TestWifisManaged(t *testing.T) {
	mock := &mockObj{}
	client := NewClient(mock)
	mock.managed = true
	res, _ := client.WifisManaged([]string{"d0", "d1"})

	if res["iface0"] != "d0" {
		t.Errorf("Expected map[iface]device not returned. Got: %v", res)
	}
	if res["iface1"] != "d1" {
		t.Errorf("Expected map[iface]device not returned. Got: %v", res)
	}
	mock.managed = false
	mock.ifaces = mock.ifaces[:0]
	res, _ = client.WifisManaged([]string{"d0", "d1"})
	if res["iface0"] == "d0" {
		t.Errorf("Expected  no result, since no ifaces are managed. Got: %v", res)
	}
}
