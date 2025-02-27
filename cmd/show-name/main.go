package main

import (
	"log"
	"os"

	"github.com/shauncampbell/golang-tplink-hs100/pkg/configuration"
	"github.com/shauncampbell/golang-tplink-hs100/pkg/hs100"
)

func main() {
	h := hs100.NewHs100("192.168.2.100", configuration.Default())

	name, err := h.GetName()
	if err != nil {
		log.Print("Error on accessing device")
		os.Exit(1)
	}

	log.Printf("Name of device: %s", name)
}
