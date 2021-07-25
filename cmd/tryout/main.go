package main

import (
	"time"

	"github.com/shauncampbell/golang-tplink-hs100/pkg/configuration"
	"github.com/shauncampbell/golang-tplink-hs100/pkg/hs100"
)

const defaultSleep = 2000 * time.Millisecond

func main() {
	h := hs100.NewHs100("localhost", configuration.Default())

	println("Name of device:")
	name, _ := h.GetName()
	println(name)

	time.Sleep(defaultSleep)

	println("Is on:")
	b, _ := h.IsOn()
	println(b)

	time.Sleep(defaultSleep)

	println("Turning on")
	_ = h.TurnOn()
	println("done")

	time.Sleep(defaultSleep)

	println("Is on:")
	b, _ = h.IsOn()
	println(b)

	time.Sleep(defaultSleep)

	println("Turning off")
	_ = h.TurnOff()
	println("done")

	time.Sleep(defaultSleep)

	println("Is on:")
	b, _ = h.IsOn()
	println(b)
}
