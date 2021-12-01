package openSdk_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGenerateSign(t *testing.T) {
	const SHORT_FORMAT_TIME = "2006-01-02 15:04:05"
	loc,_:=time.LoadLocation("Asia/Shanghai")
	var resp RequestBody
	resp.AppId = "88888888"
	resp.AppSecret = "8ddcff3a80f4189ca1c9d4d902c3c909"
	resp.Date = time.Now().In(loc).Format(SHORT_FORMAT_TIME)
	resp.Method = "analyze.tlj"
	resp.Param["content"] = "0长按复制这段文字，打开「淘→寳」即可「领取优惠券」并购买 ₰kv2YXyS89Xl£/"
	content := resp.GenerateSign()
	rp,err := http.Post("http://v3.api.haodanku.com/rest","application/json;charset=utf-8",strings.NewReader(content))
	if err != nil{
		log.Fatalln(err)
	}
	defer rp.Body.Close()
	b,err := ioutil.ReadAll(rp.Body)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(string(b))
}

func TestTmp(t *testing.T) {
	const SHORT_FORMAT_TIME = "2006-01-02 15:04:05"
	loc,_:=time.LoadLocation("Asia/Shanghai")
	var resp RequestBody
	resp.AppId = "88888888"
	resp.AppSecret = "8ddcff3a80f4189ca1c9d4d902c3c909"
	resp.Date = time.Now().In(loc).Format(SHORT_FORMAT_TIME)
	resp.Method = "analyze.tlj"
	resp.Param["content"] = "0长按复制这段文字，打开「淘→寳」即可「领取优惠券」并购买 ₰kv2YXyS89Xl£/"
	content := resp.GenerateSign()
	rp,err := http.Post("http://v3.api.haodanku.com/rest","application/json;charset=utf-8",strings.NewReader(content))
	if err != nil{
		log.Fatalln(err)
	}
	defer rp.Body.Close()
	b,err := ioutil.ReadAll(rp.Body)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println(string(b))
}
