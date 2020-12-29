package daishencha

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/godaner/ajsp/flag"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func M() {
	newHC := func() *http.Client {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			// CheckRedirect: op[0].CheckRedirect,
			// Jar:           op[0].Jar,
			Timeout: time.Duration(60) * time.Second,
		}
	}
	// read count

	// sync count
	rsc := make(chan interface{}, flag.N)
	// fetch
	go func() {
		hc := newHC()
		// loop fetch
		for ; ; {
			rs := fetch(hc)
			for _, r := range rs {
				rsc <- r
			}
			if len(rs) == 0 {
				<-time.After(time.Duration(flag.WT) * time.Second) // 等待n s
			}
		}
	}()
	// do
	for i := 0; i < int(flag.W); i++ {
		// many worker
		go func() {
			hc := newHC()
			// loop do
			for ; ; {
				do(hc, <-rsc)
			}
		}()
	}
	select {}
}

// fetch
func fetch(hc *http.Client) (rs []interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("dsc fetch : PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC , err is ", err)
		}
	}()

	req, err := http.NewRequest("POST", flag.HttpAddr+"/api/approval/wf/task/upcoming", bytes.NewBufferString(`{"page":1,"rows":`+fmt.Sprint(flag.N)+`,"orderBy":{"updateTime":"desc"},"affairName":null,"likeMap":{},"between":{},"in":{"bizStatus":["21","22","23","24","25","26","27,","28","31","33"]}}`))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", flag.Auth)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "route", Value: flag.Route, Expires: time.Now().Add(999 * time.Hour)})
	// Cookie
	//	route=6a3e9cd8159592ade356af0234099d6f
	resp, err := hc.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	log.Printf("dsc get record success , len(result) is : %+v!", len(result))
	rs = result["data"].(map[string]interface{})["records"].([]interface{})

	return rs
}

// do
func do(hc *http.Client, r interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("dsc do : PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC , err is ", err)
		}
	}()
	rm := r.(map[string]interface{})
	shardKey := rm["shardKey"].(string)
	affairId := rm["busiId"].(string)
	req, err := http.NewRequest("POST", flag.HttpAddr+"/api/approval/dth/affair/submitAffair", bytes.NewBufferString(`{"affairId":"`+affairId+`","shardKey":"`+shardKey+`","handleType":6,"auditAdvice":"1","auditAdviceInfo":"同意","base64Json":"","attIds":""}`))
	if err != nil {
		log.Printf("dsc do record err 1 , err is : %v!", err)
		return
	}
	req.Header.Set("Authorization", flag.Auth)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "route", Value: flag.Route, Expires: time.Now().Add(999 * time.Hour)})
	resp, err := hc.Do(req)
	if err != nil {
		log.Printf("dsc do record err 2 , err is : %v!", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("dsc do record err 3 , err is : %v!", err)
		return
	}
}
