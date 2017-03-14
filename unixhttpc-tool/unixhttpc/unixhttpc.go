package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func reader(r io.Reader) {
	br := bufio.NewReader(r)
	// output each line.
	line, _, err := br.ReadLine()
	for len(line) != 0 && err == nil {
		line, _, err = br.ReadLine()
	}

	if err != nil {
		log.Printf("Error reading response headers: %v\n", err)
	}

	/*// if not error, next unique line will be the json body
	line, isPrefix, err = br.ReadLine()
	body := string(line)
	for isPrefix == true && err == nil {
		line, isPrefix, err = br.ReadLine()
		body += string(line)
	}

	if err != nil {
		log.Printf("Error reading response body: %v\n", err)
	}*/

	io.Copy(os.Stdout, br)
}

func main() {
	if len(os.Args) != 3 {
		log.Println(os.Stderr, "usage:", os.Args[0], "/path.socket /uri")
		return
	}
	log.Println("Unix HTTP client requesting")

	conn, err := net.Dial("unix", os.Args[1])
	if err != nil {
		log.Println("Dial error: ", err)
		return
	}
	defer conn.Close()

	// reading here async, as should behave in a library.
	// If wanted to have the code as is, this should not be a goroutine
	// but reading sync after sending request
	go reader(conn)

	requestMsg := "GET http://unix" + os.Args[2] + " HTTP/1.1\r\n" +
		"Host: localhost\r\n" +
		"\r\n"

	_, err = conn.Write([]byte(requestMsg))
	if err != nil {
		log.Println("write error: ", err)
		return
	}

	time.Sleep(1e9)
}
