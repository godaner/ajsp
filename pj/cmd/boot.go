package main

import (
	"github.com/godaner/ajsp/pj"
)
import (
	"flag"
)

var httpAddr, sk, d string
var w, wt, s, e uint64

func main() {
	flag.StringVar(&httpAddr, "ha", "http://hcp.sczwfw.gov.cn", "http addr")
	flag.Uint64Var(&w, "w", 1, "worker")
	flag.Uint64Var(&wt, "wt", 1000, "wait time , mill sec")
	flag.StringVar(&sk, "sk", "510904", "sk")
	flag.StringVar(&d, "d", "20201230", "d")
	flag.Uint64Var(&s, "s", 0, "s")
	flag.Uint64Var(&e, "e", 0, "e")
	flag.Parse()
	if httpAddr == "" || w == 0 || wt == 0 || s == 0 || e == 0 {
		flag.PrintDefaults()
		return
	}
	go func() {
		pj.M(httpAddr, w, wt, sk, s, e, d)
	}()
	select {}
}
