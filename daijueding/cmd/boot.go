package main

import (
	"github.com/godaner/ajsp/daijueding"
	"github.com/godaner/ajsp/flag"
)

func main() {
	err := flag.Parse()
	if err != nil {
		return
	}
	go func() {
		daijueding.M()
	}()
	select {}
}
