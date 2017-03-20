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

// Client struct exposing wifi-ap operations
type Client struct {
	restClient *RestClient
}

// NewClient returns pointer to a new wifi-ap client using certain transport
func NewClient(tc TransportClient) *Client {
	return &Client{restClient: newRestClient(tc)}
}

// DefaultClient returns pointer to a new wifi-ap client using default transport
func DefaultClient() *Client {
	return &Client{restClient: defaultRestClient()}
}

func defaultServiceURI() string {
	return fmt.Sprintf("http://unix%s", filepath.Join(versionURI, configurationURI))
}

// Show shows current wifi-ap status
func (client *Client) Show() {
	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "GET", nil)
	if err != nil {
		log.Printf("wifi-ap show operation failed: %q\n", err)
		return
	}

	printMapSorted(response.Result)
}

// Enable enables wifi ap
func (client *Client) Enable() {
	params := map[string]string{"disabled": "false"}
	b, err := json.Marshal(params)
	if err != nil {
		log.Printf("wifi-ap enable operation failed when marshalling input parameters: %q\n", err)
		return
	}

	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "POST", bytes.NewReader(b))
	if err != nil {
		log.Printf("wifi-ap enable operation failed: %q\n", err)
		return
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		log.Printf("Failed to set configuration, service returned: %d (%s)\n", response.StatusCode, response.Status)
	}
}

// Disable disables wifi ap
func (client *Client) Disable() {
	params := map[string]string{"disabled": "true"}
	b, err := json.Marshal(params)
	if err != nil {
		log.Printf("wifi-ap disable operation failed when marshalling input parameters: %q\n", err)
		return
	}

	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "POST", bytes.NewReader(b))
	if err != nil {
		log.Printf("wifi-ap disable operation failed: %q\n", err)
		return
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		log.Printf("Failed to set configuration, service returned: %d (%s)\n", response.StatusCode, response.Status)
	}
}

// SetSsid sets the ssid for the wifi ap
func (client *Client) SetSsid(ssid string) {
	params := map[string]string{"wifi.ssid": ssid}
	b, err := json.Marshal(params)
	if err != nil {
		log.Printf("wifi-ap set SSID operation failed when marshalling input parameters: %q\n", err)
		return
	}

	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "POST", bytes.NewReader(b))
	if err != nil {
		log.Printf("wifi-ap set SSID operation failed: %q\n", err)
		return
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		log.Printf("Failed to set configuration, service returned: %d (%s)\n", response.StatusCode, response.Status)
	}
}

// SetPassphrase sets the credential to access the wifi ap
func (client *Client) SetPassphrase(passphrase string) {
	params := map[string]string{
		"wifi.security":            "wpa2",
		"wifi.security-passphrase": passphrase,
	}
	b, err := json.Marshal(params)
	if err != nil {
		log.Printf("wifi-ap set passphrase operation failed when marshalling input parameters: %q\n", err)
		return
	}

	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "POST", bytes.NewReader(b))
	if err != nil {
		log.Printf("wifi-ap set passphrase operation failed: %q\n", err)
		return
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		log.Printf("Failed to set configuration, service returned: %d (%s)\n", response.StatusCode, response.Status)
	}
}
