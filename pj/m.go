package pj

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

func M(httpAddr string, w, wt uint64, skey string, s, e uint64, d string) {
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
			log.Println("key is", key)
			do(hc, httpAddr, key)
			number++
			<-time.After(time.Duration(wt) * time.Millisecond)
		}()
		if end {
			panic("finish")
		}
	}
}
func do(hc *http.Client, httpAddr, key string) {
	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)

	// 取出内容类型
	contentType := w.FormDataContentType()

	// 将文件数据写入
	// h := make(textproto.MIMEHeader)
	// h.Set("Content-Disposition",
	// 	fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
	// 		"file", // 参数名为file
	// 		fname))
	// h.Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document") // 设置文件格式
	// h.Set("Content-Type", "text/plain") // 设置文件格式
	// w.WriteField("filename", fname)
	// pa, _ := w.CreatePart(h)
	// pa.Write(file_data)

	// 设置其他参数
	// w.WriteField("creatorName", uname)
	// w.WriteField("creatorId", uid)
	// w.WriteField("busiAliasCode", "apply")
	// w.WriteField("categoryCode", "apply")
	// w.WriteField("busiId", applyID)
	// w.WriteField("shardKey", shardKey)
	// w.WriteField("needToUploadPDF", "NO")
	w.WriteField("content", `{"objType":"1","objCode":"code-001","busiCode":"`+key+`","score":"5","evalSource":"1","evalChannel":"1","anonymous":"0","publish":"1","instructions":"ºÜºÃ","evaluationChildren":[{"childrenId":"4024168392809385984","childrenType":"2"},{"childrenId":"4024175170972422144","childrenType":"2"},{"childrenId":"4554175170972422144","childrenType":"2"},{"childrenId":"4024111640058523222","childrenType":"2"},{"childrenId":"4024111640058523648","childrenType":"2"},{"childrenId":"4024177111158523648","childrenType":"2"},{"childrenId":"4024177640011113648","childrenType":"2"},{"childrenId":"4024222640058523648","childrenType":"2"}],"praise":[{"paraisCongigId":"1","paraisCongigCode":"code-1"},{"paraisCongigId":"2","paraisCongigCode":"code-2"},{"paraisCongigId":"3","paraisCongigCode":"code-3"},{"paraisCongigId":"4","paraisCongigCode":"code-4"},{"paraisCongigId":"5","paraisCongigCode":"code-5"},{"paraisCongigId":"6","paraisCongigCode":"code-6"},{"paraisCongigId":"7","paraisCongigCode":"code-7"}]}`)
	w.Close()
	req, err := http.NewRequest("POST", httpAddr+"/app/api/evaluationFileHandler", b)
	if err != nil {
		log.Printf("do err 1 , err is : %v!", err)
		return
	}
	// log.Println("ctype", contentType)
	// log.Printf("b %+v", b.String())
	req.Header.Set("Content-Type", contentType)
	// req.AddCookie(&http.Cookie{Name: "topSession", Value: session, Expires: time.Now().Add(999 * time.Hour)})
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
	// result := make(map[string]interface{})
	// err = json.Unmarshal(body, &result)
	// if err != nil {
	// 	log.Printf("do err 4 , err is : %v!", err)
	// 	return
	// }
	log.Printf("do result is : %v!", string(body))
	return
}
