package wifiap

import (
	"log"
)

// Show shows current wifi-ap status
func Show() {
	result, err := DefaultRestClient().Show()
	if err != nil {
		log.Printf("wifi-ap show operation failed: %q\n", err)
		return
	}
	// TODO see if needed here to return value or simply showing that in stdout is ok
	printMapSorted(result)
}

// Enable enables wifi ap
func Enable() {
	err := DefaultRestClient().Enable()
	if err != nil {
		log.Printf("wifi-ap enable operation failed: %q\n", err)
	}
}

// Disable disables wifi ap
func Disable() {
	err := DefaultRestClient().Disable()
	if err != nil {
		log.Printf("wifi-ap disable operation failed: %q\n", err)
	}
}

// SetSsid sets the ssid for the wifi ap
func SetSsid(ssid string) {
	err := DefaultRestClient().SetSsid(ssid)
	if err != nil {
		log.Printf("wifi-ap set SSID operation failed: %q\n", err)
	}
}

// SetPassphrase sets the credential to access the wifi ap
func SetPassphrase(passphrase string) {
	err := DefaultRestClient().SetPassphrase(passphrase)
	if err != nil {
		log.Printf("wifi-ap set passphrase operation failed: %q\n", err)
	}
}
