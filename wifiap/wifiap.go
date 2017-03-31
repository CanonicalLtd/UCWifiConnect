// -*- Mode: Go; indent-tabs-mode: t -*-
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
func (client *Client) Show() (map[string]interface{}, error) {
	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "GET", nil)
	if err != nil {
		return nil, fmt.Errorf("wifi-ap show operation failed: %q", err)
	}

	return response.Result, nil
}

// Enabled checks if wifi-ap is up
func (client *Client) Enabled() (bool, error) {
	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "GET", nil)
	if err != nil {
		return false, err
	}
	fmt.Println(response.Result["disabled"])
	if response.Result["disabled"].(bool) {
		return false, nil
	}
	return true, nil
}

// Enable enables wifi ap
func (client *Client) Enable() error {
	params := map[string]string{"disabled": "false"}
	b, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("wifi-ap enable operation failed when marshalling input parameters: %q", err)
	}

	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "POST", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("wifi-ap enable operation failed: %q", err)
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		return fmt.Errorf("Failed to set configuration, service returned: %d (%s)", response.StatusCode, response.Status)
	}

	return nil
}

// Disable disables wifi ap
func (client *Client) Disable() error {
	params := map[string]string{"disabled": "true"}
	b, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("wifi-ap disable operation failed when marshalling input parameters: %q", err)
	}

	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "POST", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("wifi-ap disable operation failed: %q", err)
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		return fmt.Errorf("Failed to set configuration, service returned: %d (%s)", response.StatusCode, response.Status)
	}

	return nil
}

// SetSsid sets the ssid for the wifi ap
func (client *Client) SetSsid(ssid string) error {
	params := map[string]string{"wifi.ssid": ssid}
	b, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("wifi-ap set SSID operation failed when marshalling input parameters: %q", err)
	}

	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "POST", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("wifi-ap set SSID operation failed: %q", err)
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		return fmt.Errorf("Failed to set configuration, service returned: %d (%s)", response.StatusCode, response.Status)
	}

	return nil
}

// SetPassphrase sets the credential to access the wifi ap
func (client *Client) SetPassphrase(passphrase string) error {
	if len(passphrase) < 13 {
		return fmt.Errorf("Passphrase must be at least 13 chars in length. Please try again")
	}

	params := map[string]string{
		"wifi.security":            "wpa2",
		"wifi.security-passphrase": passphrase,
	}
	b, err := json.Marshal(params)
	if err != nil {
		return fmt.Errorf("wifi-ap set passphrase operation failed when marshalling input parameters: %q", err)
	}

	response, err := client.restClient.sendHTTPRequest(defaultServiceURI(), "POST", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("wifi-ap set passphrase operation failed: %q", err)
	}

	if response.StatusCode != http.StatusOK || response.Status != http.StatusText(http.StatusOK) {
		return fmt.Errorf("Failed to set configuration, service returned: %d (%s)", response.StatusCode, response.Status)
	}

	return nil
}
