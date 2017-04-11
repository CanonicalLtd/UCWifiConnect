package server

import (
	"log"
	"os"

	"github.com/guelfey/go.dbus"
)

const (
	sname   = "wificonnect"
	sport   = 8080
	stype   = "_http._tcp"
	sdomain = "local"
)

var stext = []string{"email=wificonnect@canonical.com", "jid=wificonnect@canonical.com", "status=avail"}

func shost() string {
	name, err := os.Hostname()
	if err != nil {
		log.Println("Could not get local hostname")
		name = ""
	}

	// myhostname.local for instance
	name = name + "." + sdomain
	return name
}

// RegisterService register wifi-connect avahi service
func RegisterService() error {
	dconn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	defer dconn.Close()

	obj := dconn.Object("org.freedesktop.Avahi", "/")
	var path dbus.ObjectPath
	obj.Call("org.freedesktop.Avahi.Server.EntryGroupNew", 0).Store(&path)

	//TODO TRACE
	log.Printf("PATH: %v\n", path)

	obj = dconn.Object("org.freedesktop.Avahi", path)

	var AAY [][]byte
	for _, s := range stext {
		AAY = append(AAY, []byte(s))
	}

	// http://www.dns-sd.org/ServiceTypes.html
	obj.Call("org.freedesktop.Avahi.EntryGroup.AddService", 0,
		int32(-1),     // avahi.IF_UNSPEC
		int32(-1),     // avahi.PROTO_UNSPEC
		uint32(0),     // flags
		sname,         // sname
		stype,         // stype
		sdomain,       // sdomain
		shost(),       // shost
		uint16(sport), // port
		AAY)           // text record
	obj.Call("org.freedesktop.Avahi.EntryGroup.Commit", 0)

	return nil
}
