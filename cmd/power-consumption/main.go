package main

import (
	"log"
	"os"

	"github.com/shauncampbell/golang-tplink-hs100/pkg/configuration"
	"github.com/shauncampbell/golang-tplink-hs100/pkg/hs100"
)

func main() {
	h := hs100.NewHs100("192.168.2.100", configuration.Default())

	p, err := h.GetCurrentPowerConsumption()
	if err != nil {
		log.Println("Error on accessing device")
		os.Exit(1)
	}

	log.Println("Current Power consumption:")
	log.Printf("Voltage: %fV", p.Voltage)
	log.Printf("Current: %fA", p.Current)
	log.Printf("Power: %fW", p.Power)
}
