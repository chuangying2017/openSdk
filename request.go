package openSdk

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
)

type RequestBody struct {
	AppId string
	Method string
	Date string
	AppSecret string
	Sign string
	Param map[string]string
}

func (r RequestBody) GenerateSign() string {
	if len(r.Param) <= 0 {
		return ""
	}
	param := make(map[string]string)
	for k,v :=range r.Param {
		param[k] = v
	}
	param["app_id"] = r.AppId
	param["method"] = r.Method
	param["date"] = r.Date
	arr := make([]string,0)
	for k:=range param {
		arr = append(arr,k)
	}
	sort.Strings(arr)
	var builder strings.Builder
	for k := range arr {
		pk := arr[k]
		builder.WriteString(pk)
		builder.WriteString(param[pk])
	}
	builder.WriteString(r.AppSecret)
	hasher := md5.New()
	hasher.Write([]byte(builder.String()))
	param["sign"] = strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
	b,err := json.Marshal(param)
	if err !=nil {
		return ""
	}
	return string(b)
}

func Tmp() {

}