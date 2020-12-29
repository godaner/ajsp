package main

import (
	"github.com/godaner/ajsp/daisouli"
	"github.com/godaner/ajsp/flag"
)

func main() {
	err := flag.Parse()
	if err != nil {
		return
	}
	go func() {
		daisouli.M()
	}()
	select {}
}
