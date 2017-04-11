package main

import (
	"bufio"
	"log"
	"os"

	"github.com/CanonicalLtd/UCWifiConnect/server"
)

func main() {

	err := server.RegisterService()
	if err != nil {
		log.Printf("Error registering avahi service: %v\n", err)
		return
	}

	// Wait for getch and exit
	bufio.NewReader(os.Stdin).ReadString('\n')

}
