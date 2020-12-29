package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var httpAddr, auth, route string
var n, w, wt uint64
var doCount, suCount uint64

// 第四个 待制证
func main() {
	flag.StringVar(&httpAddr, "ha", "http://59.225.201.162:8086", "http addr")
	flag.StringVar(&auth, "auth", "", "http Authorization")
	flag.Uint64Var(&n, "n", 50, "fetch number")
	flag.Uint64Var(&w, "w", 1, "worker")
	flag.StringVar(&route, "route", "", "route")
	flag.Uint64Var(&wt, "wt", 5, "no data wait time , sec")
	flag.Parse()
	if httpAddr == "" || auth == "" || n == 0 || w == 0 || route == "" || wt == 0 {
		flag.PrintDefaults()
		return
	}
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
	rsc := make(chan interface{}, n)
	// fetch
	go func() {
		hc := newHC()
		// loop fetch
		for ; ; {
			rs := fetch(hc, wt)
			for _, r := range rs {
				rsc <- r
			}
			if len(rs) == 0 {
				<-time.After(time.Duration(wt) * time.Second) // 等待n s
			}
		}
	}()
	// do
	for i := 0; i < int(w); i++ {
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
func fetch(hc *http.Client, wt uint64) (rs []interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("fetch : PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC , err is ", err)
		}
	}()

	// 查询
	req, err := http.NewRequest("POST", httpAddr+"/api/approval/wf/task/upcoming", bytes.NewBufferString(`{"page":1,"rows":`+fmt.Sprint(n)+`,"orderBy":{"updateTime":"desc"},"affairName":null,"likeMap":{},"between":{},"in":{"bizStatus":["11","13"]}}`))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "route", Value: route})
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
	log.Printf("get record success , result is : %+v!", result)
	rs = result["data"].(map[string]interface{})["records"].([]interface{})

	return rs
}

// do
func do(hc *http.Client, r interface{}) {
	index := atomic.AddUint64(&doCount, 1)
	defer func() {
		if err := recover(); err != nil {
			log.Println("do : PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC , err is ", err)
		}
	}()
	rm := r.(map[string]interface{})
	log.Printf("do record start , index is : %v , record is : %+v!", index, rm)
	affairId := rm["busiId"].(string)
	shardKey := rm["shardKey"].(string)
	// 获取材料
	req, err := http.NewRequest("GET", httpAddr+"/api/approval/dth/affair-material/queryAMatreialByAffairId?affairId="+affairId+"&shardKey="+shardKey, nil)
	if err != nil {
		log.Printf("do record err 1 , index is : %v , err is : %v!", index, err)
	}
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "route", Value: route})
	resp, err := hc.Do(req)
	if err != nil {
		log.Printf("do record err 2 , index is : %v , err is : %v!", index, err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("do record err 3 , index is : %v , err is : %v!", index, err)
		return
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("do record err 4 , index is : %v , err is : %v!", index, err)
		return
	}
	fail := false
	ds := result["data"].([]interface{})
	for _, d := range ds {
		// 通过材料
		dm := d.(map[string]interface{})
		materiaId := dm["id"].(string)
		req, err = http.NewRequest("GET", httpAddr+"/api/approval/dth/affair-material/auditResult?materiaId="+materiaId+"&auditType=01&auditRemark=", nil)
		if err != nil {
			log.Printf("do record err 5 , index is : %v , err is : %v!", index, err)
			fail = true
			break
		}
		req.Header.Set("Authorization", auth)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "route", Value: route})
		resp, err = hc.Do(req)
		if err != nil {
			log.Printf("do record err 6 , index is : %v , err is : %v!", index, err)
			fail = true
			break
		}
		body, err = ioutil.ReadAll(resp.Body)
		if body != nil {
			resp.Body.Close()
		}
		if err != nil {
			log.Printf("do record err 7 , index is : %v , err is : %v!", index, err)
			fail = true
			break
		}
		result := make(map[string]interface{})
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("do record err 8 , index is : %v , err is : %v!", index, err)
			fail = true
			break
		}
		if fmt.Sprint(result["success"]) != "true" {
			log.Printf("do record err 9 , index is : %v , result is : %v!", index, result)
			fail = true
			break
		}
	}
	if fail {
		log.Printf("do record err 10 , Matreial not ok , index is : %v , err is : %v!", index, err)
		return
	}
	req, err = http.NewRequest("POST", httpAddr+"/api/approval/dth/affair/submitAffair", bytes.NewBufferString(`{"affairId":"`+affairId+`","shardKey":"`+shardKey+`","handleType":1,"auditAdvice":"1","batchParts":""}`))
	if err != nil {
		log.Printf("do record err 11 , index is : %v , err is : %v!", index, err)
		return
	}
	req.Header.Set("Authorization", auth)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "route", Value: route})
	resp, err = hc.Do(req)
	if err != nil {
		log.Printf("do record err 12 , index is : %v , err is : %v!", index, err)
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("do record err 13 , index is : %v , err is : %v!", index, err)
		return
	}
	log.Printf("do record end , index is : %v , res is : %v!", index, string(body))
	atomic.AddUint64(&suCount, 1)
}
