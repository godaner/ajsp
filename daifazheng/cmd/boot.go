package main

import (
	"github.com/godaner/ajsp/daifazheng"
	"github.com/godaner/ajsp/flag"
)

func main() {

	err := flag.Parse()
	if err != nil {
		return
	}
	go func() {
		daifazheng.M()
	}()
	select {}
}
