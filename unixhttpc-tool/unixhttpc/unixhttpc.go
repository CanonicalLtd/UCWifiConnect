package main

import (
	"context"
	"fmt"
	"io"
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

	httpc := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", os.Args[1])
			},
		},
	}
	response, err := httpc.Get("http://unix" + os.Args[2])
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, response.Body)
}
