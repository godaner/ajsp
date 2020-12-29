package main

import (
	"github.com/godaner/ajsp/daishencha"
	"github.com/godaner/ajsp/flag"
)

func main() {
	err := flag.Parse()
	if err != nil {
		return
	}
	go func() {
		daishencha.M()
	}()
	select {}
}
