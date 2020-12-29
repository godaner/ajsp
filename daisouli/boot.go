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
	"time"
)

var httpAddr, auth, route string
var n, w uint64

// 第一个 待受理
func main() {
	flag.StringVar(&httpAddr, "ha", "http://59.225.201.162:8086", "http addr")
	flag.StringVar(&auth, "auth", "", "http Authorization")
	flag.Uint64Var(&n, "n", 50, "exec number")
	flag.Uint64Var(&w, "w", 1, "worker")
	flag.StringVar(&route, "route", "", "route")
	flag.Parse()
	if httpAddr == "" || auth == "" || n == 0 {
		flag.PrintDefaults()
		return
	}
	for i := 0; i < int(w); i++ {
		go func() {
			for ; ; {
				do()
			}
		}()
	}
	select {}
}

func do() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC , err is ", err)
		}
	}()
	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		// CheckRedirect: op[0].CheckRedirect,
		// Jar:           op[0].Jar,
		Timeout: time.Duration(60) * time.Second,
	}
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
	rs := result["data"].(map[string]interface{})["records"].([]interface{})
	if len(rs) == 0 {
		<-time.After(10 * time.Second) // 等待10s
	}
	for index, r := range rs {
		rm := r.(map[string]interface{})
		log.Printf("handle record start , index is : %v , record is : %+v!", index, rm)
		affairId := rm["busiId"].(string)
		shardKey := rm["shardKey"].(string)
		// 获取材料
		req, err := http.NewRequest("GET", httpAddr+"/api/approval/dth/affair-material/queryAMatreialByAffairId?affairId="+affairId+"&shardKey="+shardKey, nil)
		if err != nil {
			log.Printf("handle record err 1 , index is : %v , err is : %v!", index, err)
		}
		req.Header.Set("Authorization", auth)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "route", Value: route})
		resp, err := hc.Do(req)
		if err != nil {
			log.Printf("handle record err 2 , index is : %v , err is : %v!", index, err)
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)
		if body != nil {
			resp.Body.Close()
		}
		if err != nil {
			log.Printf("handle record err 3 , index is : %v , err is : %v!", index, err)
			continue
		}
		result := make(map[string]interface{})
		err = json.Unmarshal(body, &result)
		if err != nil {
			log.Printf("handle record err 4 , index is : %v , err is : %v!", index, err)
			continue
		}
		fail := false
		ds := result["data"].([]interface{})
		for _, d := range ds {
			// 通过材料
			dm := d.(map[string]interface{})
			materiaId := dm["id"].(string)
			req, err = http.NewRequest("GET", httpAddr+"/api/approval/dth/affair-material/auditResult?materiaId="+materiaId+"&auditType=01&auditRemark=", nil)
			if err != nil {
				log.Printf("handle record err 5 , index is : %v , err is : %v!", index, err)
				fail = true
				break
			}
			req.Header.Set("Authorization", auth)
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&http.Cookie{Name: "route", Value: route})
			resp, err = hc.Do(req)
			if err != nil {
				log.Printf("handle record err 6 , index is : %v , err is : %v!", index, err)
				fail = true
				break
			}
			body, err = ioutil.ReadAll(resp.Body)
			if body != nil {
				resp.Body.Close()
			}
			if err != nil {
				log.Printf("handle record err 7 , index is : %v , err is : %v!", index, err)
				fail = true
				break
			}
			result := make(map[string]interface{})
			err = json.Unmarshal(body, &result)
			if err != nil {
				log.Printf("handle record err 8 , index is : %v , err is : %v!", index, err)
				fail = true
				break
			}
			if fmt.Sprint(result["success"]) != "true" {
				log.Printf("handle record err 9 , index is : %v , result is : %v!", index, result)
				fail = true
				break
			}
		}
		if fail {
			log.Printf("handle record err 10 , Matreial not ok , index is : %v , err is : %v!", index, err)
			continue
		}
		req, err = http.NewRequest("POST", httpAddr+"/api/approval/dth/affair/submitAffair", bytes.NewBufferString(`{"affairId":"`+affairId+`","shardKey":"`+shardKey+`","handleType":1,"auditAdvice":"1","batchParts":""}`))
		if err != nil {
			log.Printf("handle record err 11 , index is : %v , err is : %v!", index, err)
			continue
		}
		req.Header.Set("Authorization", auth)
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "route", Value: route})
		resp, err = hc.Do(req)
		if err != nil {
			log.Printf("handle record err 12 , index is : %v , err is : %v!", index, err)
			continue
		}
		body, err = ioutil.ReadAll(resp.Body)
		if body != nil {
			resp.Body.Close()
		}
		if err != nil {
			log.Printf("handle record err 13 , index is : %v , err is : %v!", index, err)
			continue
		}
		log.Printf("handle record end , index is : %v , res is : %v!", index, string(body))
	}
}
