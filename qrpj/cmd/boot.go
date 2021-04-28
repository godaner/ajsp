package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"
)

/*
post#http://ajqwbb.suining12345.cn:10100/app/saveAppraise

formdata
satisfactionDegree=100&reportContext=&reportName=&reportTel=&code=809f8a537caf460981234975658ce420&type=emp&noSatisItem=


<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width,initial-scale=1.0,maximum-scale=1.0,user-scalable=no"/>
	<link rel="stylesheet" type="text/css" href="/mobile/css/base.css"/>
	<script language="JavaScript" type="text/javascript" src="/mobile/js/rem_fix.js"></script>
	<title>安居区机构编制二维码在线履职监督评价管理系统</title>
	<style type="text/css">
		body{margin:0;padding:0;background:url(/mobile/images/error_bg.jpg) no-repeat;height: 100vh;background-size: cover;overflow: hidden;}
	</style>
</head>
<body>
	<div class="success_container" style="text-align: center;height: 100%;">
		<img class="png" style="margin-top:2rem;display: inline-block;" src="/mobile/images/success.png"/>
		<p style="font-size: 0.4rem;color:#ffffff;">1分钟内不能重复评价！</p>
		<a href="/app/zxing/809f8a537caf460981234975658ce420" style="color:#ffffff;font-size:0.4rem;position: absolute;bottom: 0.1rem;left: calc(50% - 1rem);width: 2rem;">返回首页</a>
	</div>
</body>
</html>

*/
const (
	successStr = "^_^ 操作成功"
)

var httpAddr, typee string
var code string
var minInterval, maxInterval int64 // sec
var failExitCount int64
var d bool

func Parse() error {
	flag.StringVar(&httpAddr, "ha", "http://ajqwbb.suining12345.cn:10100/app/saveAppraise", "http addr")
	flag.StringVar(&typee, "type", "emp", "type, emp or dept")
	flag.StringVar(&code, "code", "", "code")
	flag.Int64Var(&minInterval, "mii", 60, "min interval")
	flag.Int64Var(&maxInterval, "mai", 180, "max interval")
	flag.Int64Var(&failExitCount, "fec", 3, "failure exit after n time")
	flag.BoolVar(&d, "d", false, "debug")
	flag.Parse()
	if httpAddr == "" || typee == "" || code == "" || minInterval <= 0 || maxInterval <= 0 || minInterval > maxInterval {
		flag.PrintDefaults()
		return errors.New("error")
	}
	return nil
}

func main() {
	// log.SetFlags(0)
	err := Parse()
	if err != nil {
		return
	}

	cFailExitCount := int64(0)
	c := uint64(0)
	go func() {
		for ; ; {
			func() {
				inv := minInterval
				if maxInterval != minInterval {
					source := rand.NewSource(time.Now().Unix())
					randnum := rand.New(source)
					inv = randnum.Int63n(maxInterval-minInterval+1) + minInterval
				}

				fmt.Printf("=====================Start %v====================="+fmt.Sprintln(), c)
				defer func() {
					fmt.Printf("===================== End  %v====================="+fmt.Sprintln()+fmt.Sprintln()+fmt.Sprintln(), c)
					c++
				}()
				log.Printf("Please wait a moment, this preogress will be execute in %vs......", inv)
				select {
				case <-time.After(time.Duration(inv) * time.Second):

					// satisfactionDegree=100
					// &reportContext=
					// &reportName=
					// &reportTel=
					// &code=d3b640fa502846caa53de88b2a3a2919
					// &type=dept
					// &noSatisItem=
					msg, suc := postWithFormData(http.MethodPost, httpAddr, &map[string]string{
						"satisfactionDegree": "100",
						"reportContext":      "",
						"reportName":         "",
						"reportTel":          "",
						"code":               code,
						"type":               typee,
						"noSatisItem":        "",
					})
					if !suc {
						log.Println("Execute fail!", msg, code, typee)
						cFailExitCount++
					} else {
						log.Println("Execute success!", msg, code, typee)
						cFailExitCount = 0
					}
					if cFailExitCount > failExitCount {
						log.Fatalf("Fail to appraise too much time, time is: %v!", cFailExitCount)
					}

				}
			}()
		}
	}()
	select {}
}

func postWithFormData(method, url string, postData *map[string]string) (msg string, success bool) {
	body := new(bytes.Buffer)
	w := multipart.NewWriter(body)
	for k, v := range *postData {
		w.WriteField(k, v)
	}
	w.Close()
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, _ := http.DefaultClient.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if d {
		log.Println(resp.StatusCode, string(data))
	}
	reg, _ := regexp.Compile(`<p style="font-size: 0.[0-9]*rem;color:#ffffff;">([\w\W]+)</p>`)
	res := reg.FindSubmatch(data)
	if len(res) <= 1 {
		if strings.Contains(string(data), successStr) {
			return successStr, true
		}
		log.Println("Can't find result!", resp.StatusCode, string(data))
		return "", false
	}
	if strings.Contains(string(data), successStr) {
		return successStr, true
	}
	return string(res[1]), false
}
