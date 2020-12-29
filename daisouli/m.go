package daisouli

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/godaner/ajsp/flag"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var fetcherWat sync.WaitGroup

func M() {
	newHC := func(t int64) *http.Client {
		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			// CheckRedirect: op[0].CheckRedirect,
			// Jar:           op[0].Jar,
			Timeout: time.Duration(t) * time.Second,
		}
	}
	// read count

	// sync count
	rsc := make(chan interface{}, flag.N)
	// fetch
	go func() {
		hc := newHC(60)
		// loop fetch
		for ; ; {
			rs := fetch(hc)
			if len(rs) == 0 {
				<-time.After(time.Duration(flag.WT) * time.Second) // 等待n s
				continue
			}
			fetcherWat.Add(len(rs))
			go func() {
				for _, r := range rs {
					rsc <- r
				}
			}()
			fetcherWat.Wait()
		}
	}()
	// do
	for i := 0; i < int(flag.W); i++ {
		// many worker
		go func() {
			hc := newHC(30)
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
			log.Println("dsl fetch : PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC , err is ", err)
		}
	}()

	// 查询
	req, err := http.NewRequest("POST", flag.HttpAddr+"/api/approval/wf/task/upcoming", bytes.NewBufferString(`{"page":1,"rows":`+fmt.Sprint(flag.N)+`,"orderBy":{"updateTime":"desc"},"affairName":null,"likeMap":{},"between":{},"in":{"bizStatus":["11","13"]}}`))
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
	rs = result["data"].(map[string]interface{})["records"].([]interface{})
	log.Printf("dsl get record success , len(result) is : %+v!", len(rs))

	return rs
}

// do
func do(hc *http.Client, r interface{}) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("dsl do : PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC , err is ", err)
		}
	}()
	defer func() {
		fetcherWat.Done()
	}()
	rm := r.(map[string]interface{})
	affairId := rm["busiId"].(string)
	shardKey := rm["shardKey"].(string)
	// 获取材料
	req, err := http.NewRequest("GET", flag.HttpAddr+"/api/approval/dth/affair-material/queryAMatreialByAffairId?affairId="+affairId+"&shardKey="+shardKey, nil)
	if err != nil {
		log.Printf("do record err 1 , err is : %v!", err)
	}
	req.Header.Set("Authorization", flag.Auth)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "route", Value: flag.Route, Expires: time.Now().Add(999 * time.Hour)})
	resp, err := hc.Do(req)
	if err != nil {
		log.Printf("dsl do record err 2 , err is : %v!", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("dsl do record err 3 , err is : %v!", err)
		return
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("dsl do record err 4 , err is : %v!", err)
		return
	}
	fail := false
	ds := result["data"].([]interface{})
	for _, d := range ds {
		// 通过材料
		dm := d.(map[string]interface{})
		materiaId := dm["id"].(string)
		req, err = http.NewRequest("GET", flag.HttpAddr+"/api/approval/dth/affair-material/auditResult?materiaId="+materiaId+"&auditType=01&auditRemark=", nil)
		if err != nil {
			log.Printf("dsl do record err 5 , err is : %v!", err)
			fail = true
			break
		}
		req.Header.Set("Authorization", flag.Auth)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "route", Value: flag.Route, Expires: time.Now().Add(999 * time.Hour)})
		resp, err = hc.Do(req)
		if err != nil {
			log.Printf("dsl do record err 6 , err is : %v!", err)
			fail = true
			break
		}
		body, err = ioutil.ReadAll(resp.Body)
		if body != nil {
			resp.Body.Close()
		}
		if err != nil {
			log.Printf("dsl do record err 7 , err is : %v!", err)
			fail = true
			break
		}
		result := make(map[string]interface{})
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("dsl do record err 8 , err is : %v!", err)
			fail = true
			break
		}
		if fmt.Sprint(result["success"]) != "true" {
			log.Printf("dsl do record err 9 , result is : %v!", result)
			fail = true
			break
		}
	}
	if fail {
		log.Printf("dsl do record err 10 , Matreial not ok , err is : %v!", err)
		return
	}
	req, err = http.NewRequest("POST", flag.HttpAddr+"/api/approval/dth/affair/submitAffair", bytes.NewBufferString(`{"affairId":"`+affairId+`","shardKey":"`+shardKey+`","handleType":1,"auditAdvice":"1","batchParts":""}`))
	if err != nil {
		log.Printf("dsl do record err 11 , err is : %v!", err)
		return
	}
	req.Header.Set("Authorization", flag.Auth)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "route", Value: flag.Route, Expires: time.Now().Add(999 * time.Hour)})
	resp, err = hc.Do(req)
	if err != nil {
		log.Printf("dsl do record err 12 , err is : %v!", err)
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("dsl do record err 13 , err is : %v!", err)
		return
	}
}
