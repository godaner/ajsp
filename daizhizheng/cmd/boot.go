package main

import (
	"github.com/godaner/ajsp/daizhizheng"
	"github.com/godaner/ajsp/flag"
)

func main() {
	err := flag.Parse()
	if err != nil {
		return
	}
	go func() {
		daizhizheng.M()
	}()
	select {}
}
