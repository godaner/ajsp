package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"regexp"
	"sync/atomic"
	"time"
)

var httpAddr, session string
var projectID, f, fname, uname, uid, shardKey, mcode, mname, id, eventCode, eventName, areaCode, depName, depCode string
var w, wt uint64
var doCount, suCount uint64

// 提单
func main() {
	// tmp
	flag.StringVar(&depName, "depName", "遂宁市安居区行政审批局", "dep name")
	flag.StringVar(&depCode, "depCode", "3938932771577843712", "dep code")
	flag.StringVar(&areaCode, "areaCode", "510904000000", "area code")
	flag.StringVar(&eventCode, "eventCode", "511A0004200014-510904000000-000-769975683-3-00", "eventCode/选项的itemcode")
	flag.StringVar(&eventName, "eventName", "放射诊疗许可证遗失补办（县级）", "eventCode/选项的itemname")
	flag.StringVar(&id, "id", "100", "id")
	flag.StringVar(&mcode, "mcode", "2020007484", "materialCode")
	flag.StringVar(&mname, "mname", "遗失补办申请", "materialName")
	flag.StringVar(&shardKey, "skey", "5109", "shardKey")
	flag.StringVar(&uid, "uid", "2007783531106795520", "userid")
	flag.StringVar(&uname, "uname", "刘惠莹", "username")
	flag.StringVar(&fname, "fname", "lhy测试件.docx", "filename")
	flag.StringVar(&f, "f", "./lhy测试件.docx", "test file")
	flag.StringVar(&projectID, "pid", "4101290569786753024", "project id")
	flag.StringVar(&session, "session", "", "session")
	// const
	flag.StringVar(&httpAddr, "ha", "http://zxbl.sczwfw.gov.cn/app", "http addr")
	flag.Uint64Var(&w, "w", 1, "worker")
	flag.Uint64Var(&wt, "wt", 1, "wait time , sec")
	flag.Parse()
	if httpAddr == "" || session == "" || w == 0 || wt == 0 {
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
	// do
	for i := 0; i < int(w); i++ {
		// many worker
		go func() {
			hc := newHC()
			// loop do
			for ; ; {
				do(hc)
			}
		}()
	}
	select {}
}

// do
func do(hc *http.Client) {
	index := atomic.AddUint64(&doCount, 1)
	defer func() {
		<-time.After(time.Duration(wt) * time.Second)
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Println("do : PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC PANIC , err is ", err)
		}
	}()
	applyID := getBusidApplyID(hc)
	if applyID == "" {
		return
	}
	fid := uploadFile(hc, applyID)
	if fid == "" {
		return
	}
	succ := applyFile(hc, applyID, fid)
	if succ == "" {
		return
	}
	succ = submit(hc, applyID, fid)
	if succ == "" {
		return
	}
	atomic.AddUint64(&suCount, 1)
	log.Printf("handle success , index is : %v!", index)
}

func submit(hc *http.Client, applyID string, fid string) string {
	req, err := http.NewRequest("POST", httpAddr+"/zwOnline/apply/createAndSubmitSubjectApply", bytes.NewBufferString(`id=`+applyID+`&status=0&payStatus=0&expressStatus=0&implListId=`+projectID+`&eventCode=`+eventCode+`&eventName=`+eventName+`&deptCode=`+depCode+`&deptName=`+depName+`&esheetOutValue={"cardType":"111","contactCardType":"111","applicantName":"刘惠莹","applicantSex":"2","applicantIdCard":"360722199509065421","applicantPhone":"18582485421","applicantEmail":"","applicantAddress":"","contactName":"刘惠莹","contactSex":"2","contactIdCard":"360722199509065421","contactPhone":"18582485421"}&shardKey=`+shardKey+`&takingWay=2&areaCode=`+areaCode+`&conditionID=&remark=fromAcceptHandle`))
	if err != nil {
		log.Printf("submit err 1 , err is : %v!", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.AddCookie(&http.Cookie{Name: "topSession", Value: session, Expires: time.Now().Add(999 * time.Hour)})
	resp, err := hc.Do(req)
	if err != nil {
		log.Printf("submit err 2 , err is : %v!", err)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("submit err 3 , err is : %v!", err)
		return ""
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("submit err 4 , err is : %v!", err)
		return ""
	}
	if fmt.Sprint(result["result"]) == "true" {
		log.Printf("submit err success , result is : %v!", string(body))
		return "true"
	}
	log.Printf("submit err 5 , err is : %v!", errors.New(string(body)))
	return ""
}

func applyFile(hc *http.Client, applyID string, fid string) string {
	req, err := http.NewRequest("POST", httpAddr+"/presonServices/createApplyFile", bytes.NewBufferString("id="+id+"&applyId="+applyID+"&materialCode="+mcode+"&valueString="+fid+"&materialItemName="+mname+"&materialId=0&fileName="+fname+"&fileType=docx&downUrl=/app/attachment/download?id="+fid+"&shardKey="+shardKey+"&ifEsheet=0"))
	if err != nil {
		log.Printf("applyFile err 1 , err is : %v!", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.AddCookie(&http.Cookie{Name: "topSession", Value: session, Expires: time.Now().Add(999 * time.Hour)})
	resp, err := hc.Do(req)
	if err != nil {
		log.Printf("applyFile err 2 , err is : %v!", err)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("applyFile err 3 , err is : %v!", err)
		return ""
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("applyFile err 4 , err is : %v!", err)
		return ""
	}
	if fmt.Sprint(result["result"]) == "true" {
		log.Printf("applyFile success result is : %v!", string(body))
		return "true"
	}
	log.Printf("applyFile err 5 , err is : %v!", errors.New(string(body)))
	return ""
}

func uploadFile(hc *http.Client, applyID string) string {

	// bodyBuf := &bytes.Buffer{}
	// bodyWriter := multipart.NewWriter(bodyBuf)
	//
	// // "file" 为接收时定义的参数名
	// fileWriter, err := bodyWriter.CreateFormFile(fname, f)
	// if err != nil {
	// 	log.Printf("uploadFile err 4 , err is : %v!", err)
	// 	return ""
	// }
	//
	// // 打开文件
	// fh, err := os.Open(f)
	// if err != nil {
	// 	log.Printf("uploadFile err 5 , err is : %v!", err)
	// 	return ""
	// }
	// defer fh.Close()
	//
	// // iocopy
	// _, err = io.Copy(fileWriter, fh)
	// if err != nil {
	// 	return ""
	// }
	//
	// contentType := bodyWriter.FormDataContentType()

	// 读出文本文件数据
	file_data, _ := ioutil.ReadFile(f)

	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)

	// 取出内容类型
	contentType := w.FormDataContentType()

	// 将文件数据写入
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			"file", // 参数名为file
			fname))
	h.Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document") // 设置文件格式
	// h.Set("Content-Type", "text/plain") // 设置文件格式
	w.WriteField("filename", fname)
	pa, _ := w.CreatePart(h)
	pa.Write(file_data)

	// 设置其他参数
	w.WriteField("creatorName", uname)
	w.WriteField("creatorId", uid)
	w.WriteField("busiAliasCode", "apply")
	w.WriteField("categoryCode", "apply")
	w.WriteField("busiId", applyID)
	w.WriteField("shardKey", shardKey)
	w.WriteField("needToUploadPDF", "NO")
	w.WriteField("fileType", "")
	w.Close()
	req, err := http.NewRequest("POST", httpAddr+"/attachments/uploadFile", b)
	if err != nil {
		log.Printf("uploadFile err 1 , err is : %v!", err)
		return ""
	}
	// log.Println("ctype", contentType)
	// log.Printf("b %+v", b.String())
	req.Header.Set("Content-Type", contentType)
	req.AddCookie(&http.Cookie{Name: "topSession", Value: session, Expires: time.Now().Add(999 * time.Hour)})
	resp, err := hc.Do(req)
	if err != nil {
		log.Printf("uploadFile err 2 , err is : %v!", err)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("uploadFile err 3 , err is : %v!", err)
		return ""
	}
	result := make(map[string]interface{})
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("uploadFile err 4 , err is : %v!", err)
		return ""
	}
	fid := fmt.Sprint(result["id"])
	log.Printf("uploadFile success fid is : %v!", fid)
	return fid
}
func getBusidApplyID(hc *http.Client) string {
	req, err := http.NewRequest("GET", httpAddr+"/presonServices/netApply/"+projectID+"/0", nil)
	if err != nil {
		log.Printf("getBusidApplyID err 1 , err is : %v!", err)
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "topSession", Value: session, Expires: time.Now().Add(999 * time.Hour)})
	resp, err := hc.Do(req)
	if err != nil {
		log.Printf("getBusidApplyID err 2 , err is : %v!", err)
		return ""
	}
	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		resp.Body.Close()
	}
	if err != nil {
		log.Printf("getBusidApplyID err 3 , err is : %v!", err)
		return ""
	}
	comp, err := regexp.Compile(`\{seriesNumber:\'([0-9]+)\'\}`)
	if err != nil {
		log.Printf("getBusidApplyID err 4 , err is : %v!", err)
		return ""
	}
	subs := comp.FindSubmatch(body)
	if len(subs) <= 1 {
		log.Printf("getBusidApplyID err 4 , err is : %v!", errors.New("no applyid"))
		return ""
	}
	applyID := string(subs[1])
	log.Printf("getBusidApplyID success applyID is : %v!", applyID)
	return applyID
}
