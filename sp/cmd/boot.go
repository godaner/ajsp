package main

import (
	"github.com/godaner/ajsp/daifazheng"
	"github.com/godaner/ajsp/daijueding"
	"github.com/godaner/ajsp/daishencha"
	"github.com/godaner/ajsp/daisouli"
	"github.com/godaner/ajsp/daizhizheng"
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
	go func() {
		daijueding.M()
	}()
	go func() {
		daishencha.M()
	}()
	go func() {
		daisouli.M()
	}()
	go func() {
		daizhizheng.M()
	}()
	select {}
}
