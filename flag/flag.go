package flag

import (
	"errors"
	"flag"
)

var HttpAddr, Auth, Route string
var N, W, WT uint64

func Parse() error {
	flag.StringVar(&HttpAddr, "ha", "http://59.225.201.162:8086", "http addr")
	flag.StringVar(&Auth, "auth", "", "http Authorization")
	flag.Uint64Var(&N, "n", 50, "fetch number")
	flag.Uint64Var(&W, "w", 1, "worker")
	flag.StringVar(&Route, "route", "", "route")
	flag.Uint64Var(&WT, "wt", 5, "no data wait time , sec")
	flag.Parse()
	if HttpAddr == "" || Auth == "" || N == 0 || W == 0 || Route == "" || WT == 0 {
		flag.PrintDefaults()
		return errors.New("error")
	}
	return nil
}
