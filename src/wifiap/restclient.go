package wifiap

import (
	"bufio"
	"log"
	"net"
)

const socketPath = "/var/snap/wifi-ap/current/sockets/control"
const versionPath = "/v1"
const configurationPath = "/configuration"

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

func (restClient *RestClient) newConn() error {
	var err error
	restClient.conn, err = net.Dial("unix", restClient.SocketPath)
	if err != nil {
		log.Printf("Dial error: %v\n", err)
	}
	return err
}

func (restClient *RestClient) sendRequest(request string) error {
	_, err := restClient.conn.Write([]byte(request))
	if err != nil {
		log.Printf("Write error: %v\n", err)
		return err
	}

	return nil
}

// Show renders
func (restClient *RestClient) Show() (string, error) {

	err := restClient.newConn()
	if err != nil {
		return "", err
	}
	defer restClient.conn.Close()

	requestMsg := "GET http://unix" + versionPath + configurationPath + " HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"\r\n"

	err = restClient.sendRequest(requestMsg)
	if err != nil {
		log.Printf("Error sending request: %v\n", err)
		return "", err
	}

	return restClient.getResponse()
}

func (restClient *RestClient) getResponse() (string, error) {
	br := bufio.NewReader(restClient.conn)
	// read header lines and do nothing with them
	line, _, err := br.ReadLine()
	for len(line) != 0 && err == nil {
		line, _, err = br.ReadLine()
	}

	if err != nil {
		log.Printf("Error reading response headers: %v\n", err)
		return "", err
	}

	// next line is the body
	// TODO VERIFY in all cases the body comes in only one line
	body, err := br.ReadString('\n')
	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
		return "", err
	}

	return body, nil
}
