package pj

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

func M(httpAddr string, w, wt uint64, skey string, s, e uint64, d, c string) {
	hc := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		// CheckRedirect: op[0].CheckRedirect,
		// Jar:           op[0].Jar,
		Timeout: time.Duration(60) * time.Second,
	}
	// nows := time.Now().Local().Format("20060102")
	// log.Println("nows is", nows)
	number := s
	end := false
	cc := strings.Split(c, ",")
	for ; ; {
		func() {
			defer func() {
				if err := recover(); err != nil {
					log.Println("PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC , err is ", err)
				}
			}()
			log.Println("number is", number)
			if number > e {
				end = true
				return
			}
			// 510904-20201230-001638
			key := strings.Replace(fmt.Sprintf("%v-%v-%6d", skey, d, number), " ", "0", -1)
			do(hc, httpAddr, key, cc)
			fmt.Println("=======================================================")
			number++
			t := random(int64(wt))
			if t == -1 {
				t = int64(wt)
			}
			tt := time.Duration(t) * time.Millisecond
			log.Println("wait ", tt, ", for next ......")
			<-time.After(tt)
		}()
		if end {
			panic("finish")
		}
	}
}
func do(hc *http.Client, httpAddr, key string, cc []string) {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	i := random(int64(len(cc)))
	if i == -1 {
		i = 0
	}
	ins := cc[i]
	log.Println("key is", key, ", instructions is", ins)
	// 取出内容类型
	contentType := w.FormDataContentType()
	w.WriteField("content", `{"objType":"1","objCode":"code-001","busiCode":"`+key+`","score":"5","evalSource":"1","evalChannel":"1","anonymous":"0","publish":"1","instructions":"`+ins+`","evaluationChildren":[{"childrenId":"4024168392809385984","childrenType":"2"},{"childrenId":"4024175170972422144","childrenType":"2"},{"childrenId":"4554175170972422144","childrenType":"2"},{"childrenId":"4024111640058523222","childrenType":"2"},{"childrenId":"4024111640058523648","childrenType":"2"},{"childrenId":"4024177111158523648","childrenType":"2"},{"childrenId":"4024177640011113648","childrenType":"2"},{"childrenId":"4024222640058523648","childrenType":"2"}],"praise":[{"paraisCongigId":"1","paraisCongigCode":"code-1"},{"paraisCongigId":"2","paraisCongigCode":"code-2"},{"paraisCongigId":"3","paraisCongigCode":"code-3"},{"paraisCongigId":"4","paraisCongigCode":"code-4"},{"paraisCongigId":"5","paraisCongigCode":"code-5"},{"paraisCongigId":"6","paraisCongigCode":"code-6"},{"paraisCongigId":"7","paraisCongigCode":"code-7"}]}`)
	w.Close()
	req, err := http.NewRequest("POST", httpAddr+"/app/api/evaluationFileHandler", b)
	if err != nil {
		log.Printf("do err 1 , err is : %v!", err)
		return
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := hc.Do(req)
	if err != nil {
		log.Printf("do err 2 , err is : %v!", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("do err 3 , err is : %v!", err)
		return
	}
	log.Printf("do result is : %v!", string(body))
	return
}

func random(max int64) int64 {
	var i *big.Int
	i, err := rand.Int(rand.Reader, new(big.Int).SetInt64(max))
	if err != nil {
		return -1
	}
	return i.Int64()
}
