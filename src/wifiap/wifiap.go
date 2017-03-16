//
// Copyright (C) 2017 Canonical Ltd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License version 3 as
// published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package wifiap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
)

// Show shows current wifi-ap status
func Show() {
	uri := fmt.Sprintf("http://unix%s", filepath.Join(versionURI, configurationURI))
	response, err := defaultRestClient().sendHTTPRequest(uri, "GET", nil)
	if err != nil {
		log.Printf("wifi-ap show operation failed: %q\n", err)
		return
	}

	printMapSorted(response.Result)
}

// Enable enables wifi ap
func Enable() {
	params := map[string]string{"disabled": "false"}
	b, err := json.Marshal(params)
	if err != nil {
		log.Printf("wifi-ap enable operation failed when marshalling input parameters: %q\n", err)
		return
	}

	uri := fmt.Sprintf("http://unix%s", filepath.Join(versionURI, configurationURI))
	response, err := defaultRestClient().sendHTTPRequest(uri, "POST", bytes.NewReader(b))
	if err != nil {
		log.Printf("wifi-ap enable operation failed: %q\n", err)
		return
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		log.Printf("Failed to set configuration, service returned: %d (%s)\n", response.StatusCode, response.Status)
	}
}

// Disable disables wifi ap
func Disable() {
	params := map[string]string{"disabled": "true"}
	b, err := json.Marshal(params)
	if err != nil {
		log.Printf("wifi-ap disable operation failed when marshalling input parameters: %q\n", err)
		return
	}

	uri := fmt.Sprintf("http://unix%s", filepath.Join(versionURI, configurationURI))
	response, err := defaultRestClient().sendHTTPRequest(uri, "POST", bytes.NewReader(b))
	if err != nil {
		log.Printf("wifi-ap disable operation failed: %q\n", err)
		return
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		log.Printf("Failed to set configuration, service returned: %d (%s)\n", response.StatusCode, response.Status)
	}
}

// SetSsid sets the ssid for the wifi ap
func SetSsid(ssid string) {
	params := map[string]string{"wifi.ssid": ssid}
	b, err := json.Marshal(params)
	if err != nil {
		log.Printf("wifi-ap set SSID operation failed when marshalling input parameters: %q\n", err)
		return
	}

	uri := fmt.Sprintf("http://unix%s", filepath.Join(versionURI, configurationURI))
	response, err := defaultRestClient().sendHTTPRequest(uri, "POST", bytes.NewReader(b))
	if err != nil {
		log.Printf("wifi-ap set SSID operation failed: %q\n", err)
		return
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		log.Printf("Failed to set configuration, service returned: %d (%s)\n", response.StatusCode, response.Status)
	}
}

// SetPassphrase sets the credential to access the wifi ap
func SetPassphrase(passphrase string) {
	params := map[string]string{
		"wifi.security":            "wpa2",
		"wifi.security-passphrase": passphrase,
	}
	b, err := json.Marshal(params)
	if err != nil {
		log.Printf("wifi-ap set passphrase operation failed when marshalling input parameters: %q\n", err)
		return
	}

	uri := fmt.Sprintf("http://unix%s", filepath.Join(versionURI, configurationURI))
	response, err := defaultRestClient().sendHTTPRequest(uri, "POST", bytes.NewReader(b))
	if err != nil {
		log.Printf("wifi-ap set passphrase operation failed: %q\n", err)
		return
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		log.Printf("Failed to set configuration, service returned: %d (%s)\n", response.StatusCode, response.Status)
	}
}
