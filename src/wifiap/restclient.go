package wifiap

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"path/filepath"
)

const (
	socketPath       = "/var/snap/wifi-ap/current/sockets/control"
	versionURI       = "/v1"
	configurationURI = "/configuration"
)

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
func NewRestClient(socketPath string) *RestClient {
	return &RestClient{SocketPath: socketPath}
}

// DefaultRestClient created a RestClient object pointing to default socket path
func DefaultRestClient() *RestClient {
	return NewRestClient(socketPath) // FIXME this would be better like os.Getenv("SNAP_COMMON") + "/sockets/control", but that's not working
}

func (restClient *RestClient) unixDialer(_, _ string) (net.Conn, error) {
	return net.Dial("unix", restClient.SocketPath)
}

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

// Show renders
func (restClient *RestClient) Show() (map[string]interface{}, error) {
	uri := fmt.Sprintf("http://unix%s", filepath.Join(versionURI, configurationURI))
	response, err := restClient.sendHTTPRequest(uri, "GET", nil)
	if err != nil {
		return nil, err
	}

	return response.Result, nil
}
