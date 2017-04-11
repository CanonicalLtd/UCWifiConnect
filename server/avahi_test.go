package server

import (
	"log"
	"testing"

	"github.com/godbus/dbus"
)

func lookup() error {
	dconn, err := dbus.SystemBus()
	if err != nil {
		return err
	}

	defer dconn.Close()

	log.Printf("Sysdbus: %+v", dconn)

	avahiObject := dconn.Object("org.freedesktop.Avahi", "/")

	log.Printf("Avahi object: %+v", avahiObject)

	// resolve service
	avahiObject.Call("org.freedesktop.Avahi.Server.ResolveService", 0,
		int32(-1), // avahi.IF_UNSPEC
		int32(-1), // avahi.PROTO_UNSPEC
		uint32(0), // flags
	)

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

	return nil
}

func TestLookup(t *testing.T) {
	lookup()
}
