/*package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.socket /uri")
		return
	}
	fmt.Println("Unix HTTP client")

	req, err := http.NewRequest("GET", os.Args[3], nil)
	if err != nil {
		log.Fatal("New request error", err)
	}

	var resp *http.Response

	httpc := &http.Client{
		Transport: &http.Transport{
			Dial: net.Dial("unix", os.Args[1]),
		},
	}
	resp, err = httpc.Do(req)
	if err != nil {
		log.Fatal("Do request error", err)
	}

	defer resp.Body.Close()

	io.Copy(os.Stdout, resp.Body)
}*/

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		println("Client got:", string(buf[0:n]))

		//io.Copy(os.Stdout, response.Body)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.socket /uri")
		return
	}
	fmt.Println("Unix HTTP client.")

	c, err := net.Dial("unix", os.Args[1])
	if err != nil {
		fmt.Println("Dial error: ", err)
		return
	}
	defer c.Close()

	go reader(c)

	requestMsg := "GET http://unix" + os.Args[2] + " HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"\r\n"

	_, err = c.Write([]byte(requestMsg))
	if err != nil {
		fmt.Println("write error: ", err)
		return
	}
}
