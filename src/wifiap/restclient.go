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
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

const (
	versionURI       = "/v1"
	configurationURI = "/configuration"
)

var socketPath = os.Getenv("SNAP_COMMON") + "/sockets/control"

type serviceResponse struct {
	Result     map[string]interface{} `json:"result"`
	Status     string                 `json:"status"`
	StatusCode int                    `json:"status-code"`
	Type       string                 `json:"type"`
}

// RestClient defines client for rest api exposed by a unix socket
type RestClient struct {
	SocketPath string
	conn       net.Conn
}

// NewRestClient creates a RestClient object pointing to socket path set as parameter
func newRestClient(socketPath string) *RestClient {
	return &RestClient{SocketPath: socketPath}
}

// DefaultRestClient created a RestClient object pointing to default socket path
func defaultRestClient() *RestClient {
	return newRestClient(socketPath) // FIXME this would be better like os.Getenv("SNAP_COMMON") + "/sockets/control", but that's not working
}

func (restClient *RestClient) unixDialer(_, _ string) (net.Conn, error) {
	return net.Dial("unix", restClient.SocketPath)
}

// SendHTTPRequest sends a HTTP request to certain URI, using certain method and providing json parameters if needed
func (restClient *RestClient) sendHTTPRequest(uri string, method string, body io.Reader) (*serviceResponse, error) {
	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Dial: restClient.unixDialer,
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	realResponse := &serviceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&realResponse); err != nil {
		return nil, err
	}

	if realResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed: %s", realResponse.Result["message"])
	}

	return realResponse, nil
}
