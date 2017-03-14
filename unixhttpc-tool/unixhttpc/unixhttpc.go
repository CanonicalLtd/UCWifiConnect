package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage:", os.Args[0], "/path.socket /uri")
		return
	}
	fmt.Println("Unix HTTP client")

	httpc := http.Client{
		Timeout: time.Second * 5,
	}

	response, err := httpc.Get("http://unix" + os.Args[2])
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, response.Body)
}
