package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
)

type TaoBaoResqBody struct {
	method      string
	app_key     string
	session     string
	timestamp   string
	v           string
	sign_method string
	sign        string
	format      string
	simplify    string
}

var RwTpwdQueryRequest RepWirelessShareTpwdQueryRequest

func init() {
	RwTpwdQueryRequest.v = apiVersion
	RwTpwdQueryRequest.format = format
	RwTpwdQueryRequest.app_key = appKey
	RwTpwdQueryRequest.sign_method = signMethod
	RwTpwdQueryRequest.simplify = "true"
}

const (
	gatewayUrl      = "http://gw.api.taobao.com/router/rest"
	format          = "json"
	checkRequest    = true
	signMethod      = "md5"
	apiVersion      = "2.0"
	sdkVersion      = "top-sdk-php-20180326"
	appKey          = "23565334"
	secretKey       = "e6f6610a14cb95c27dacec0938181be5"
	ShortFormatTime = "2006-01-02 15:04:05"
)

type RepWirelessShareTpwdQueryRequest struct {
	TaoBaoResqBody
	password_content string
}

func (r *TaoBaoResqBody) GenerateSign(mp map[string]string) map[string]string {
	var builder strings.Builder
	loc, _ := time.LoadLocation("Asia/Shanghai")
	builder.WriteString(secretKey)
	mp["method"] = r.method
	mp["simplify"] = r.simplify
	mp["v"] = r.v
	mp["format"] = r.format
	mp["sign_method"] = r.sign_method
	mp["app_key"] = r.app_key
	mp["timestamp"] = time.Now().In(loc).Format(ShortFormatTime)
	arr := make([]string, len(mp))
	for k := range mp {
		arr = append(arr, k)
	}
	sort.Strings(arr)
	for k := range arr {
		builder.WriteString(arr[k])
		builder.WriteString(mp[arr[k]])
	}
	builder.WriteString(secretKey)
	hasher := md5.New()
	hasher.Write([]byte(builder.String()))
	mp["sign"] = strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
	return mp
}

func (r *RepWirelessShareTpwdQueryRequest) Request() ([]byte, error) {
	r.method = "taobao.wireless.share.tpwd.query"
	respMap := r.GenerateSign(map[string]string{
		"password_content": r.password_content,
	})
	DataUrlVal := url.Values{}
	for key, val := range respMap {
		DataUrlVal.Add(key, val)
	}
	rp, err := http.Post(gatewayUrl, "application/x-www-form-urlencoded; charset=UTF-8", strings.NewReader(DataUrlVal.Encode()))

	if err != nil {
		return nil, err
	}

	defer rp.Body.Close()

	b, err := ioutil.ReadAll(rp.Body)

	if err != nil {
		return nil, err
	}

	return b, nil
}

type rpJson struct {
	Code uint16 `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

func defaultHttp(w http.ResponseWriter, r *http.Request) {
	path, httpMethod := r.URL.Path, r.Method
	var err error
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	rp := rpJson{
		Code: 500,
		Msg:  "failure",
		Data: "",
	}
	if httpMethod != "GET" {
		rp.Msg = "请求方法不被允许"
		rp.Code = 401
		b1, _ := json.Marshal(rp)
		w.Write(b1)
		return
	}
	if path != "/taowordparse" {
		rp.Msg = "请求路径不存在"
		rp.Code = 404
		b1, _ := json.Marshal(rp)
		w.Write(b1)
		return
	}
	taoword := r.URL.Query().Get("taoword")
	taoword, err = url.QueryUnescape(taoword)
	if taoword == "" || err != nil {
		rp.Msg = "请求参数有误"
		rp.Code = 1001
		b1, _ := json.Marshal(rp)
		w.Write(b1)
		return
	}
	RwTpwdQueryRequest.password_content = taoword
	resp, err := RwTpwdQueryRequest.Request()
	if err != nil {
		rp.Msg = err.Error()
		rp.Code = 502
		b1, _ := json.Marshal(rp)
		w.Write(b1)
	} else {
		w.Write(resp)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := http.Server{
		Addr:    ":8099",
		Handler: http.TimeoutHandler(http.HandlerFunc(defaultHttp), 2*time.Second, "Timeout!!!"),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Wait for OS exit signal
	<-exit
	srv.Shutdown(ctx)
	log.Println("Got exit signal")
	os.Exit(0)
}
