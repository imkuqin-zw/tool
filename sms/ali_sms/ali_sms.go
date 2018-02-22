package ali_sms

import (
	"github.com/google/uuid"
	"time"
	"encoding/base64"
	"net/url"
	"strings"
	"crypto/sha1"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"errors"
	"crypto/hmac"
	"github.com/imkuqin-zw/tool/order"
)

type SendResponse struct {
	RequestId 	string
	Code		string
	Message 	string
	BizId		string
}

var (
	signName = ""
	verifyTemplateCode = ""
	accessKeyId = ""
	accessKeySecret = ""
)

func request(accessKeyId, accessKeySecret, domain string, params map[string]string, security bool) (body []byte, err error) {
	kv := order.NewKVMap(6)
	kv.Set("SignatureMethod", "HMAC-SHA1")
	kv.Set("SignatureNonce", uuid.New().String())
	kv.Set("SignatureVersion", "1.0")
	kv.Set("AccessKeyId", accessKeyId)
	kv.Set("Format", "JSON")
	kv.Set("Timestamp", time.Now().UTC().Format(time.RFC3339))
	for k, v := range params {
		kv.Set(k, v)
	}
	key, val := kv.SortAsc().GetAllKV()
	var sortedQueryStringTmp = ""
	for i, k := range key {
		sortedQueryStringTmp += "&" + encode(k) + "=" + encode(val[i].(string))
	}
	signature, err := signatureMethod("GET&%2F&" + encode(sortedQueryStringTmp[1:]), accessKeySecret)
	if err != nil {
		return
	}
	var url string
	if security {
		url = "https"
	} else {
		url = "http"
	}
	url += "://" + domain + "/?Signature=" + signature + sortedQueryStringTmp
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func signatureMethod(src, accessSecret string) (signature string, err error) {
	mac := hmac.New(sha1.New, []byte(accessSecret + "&"))
	_, err = mac.Write([]byte(src))
	if err != nil {
		return
	}
	signature = encode(base64.StdEncoding.EncodeToString(mac.Sum(nil)))
	return
}

func encode(src string) (out string) {
	out = url.QueryEscape(src)
	out = strings.Replace(out, "+", "%20", -1)
	out = strings.Replace(out, "*", "%2A", -1)
	out = strings.Replace(out, "%7E", "~", -1)
	return
}

func SendVerifySms(phone string, code string) (err error) {
	params := make(map[string]string)
	params["PhoneNumbers"] = phone
	params["SignName"] = signName
	params["TemplateCode"] = verifyTemplateCode
	templateParam, err := json.Marshal(map[string]string{"code": code})
	if err != nil {
		return
	}
	params["TemplateParam"] = string(templateParam)
	params["RegionId"] = string("cn-hangzhou")
	params["Action"] = string("SendSms")
	params["Version"] = string("2017-05-25")
	body, err := request(accessKeyId, accessKeySecret, "dysmsapi.aliyuncs.com", params, false)
	if err != nil {
		return
	}
	var resp SendResponse
	if err = json.Unmarshal(body, &resp); err != nil {
		return
	}
	if resp.Code != "OK" {
		err = errors.New(resp.Message)
	}
	return
}

func init() {
	//从配置文件等地方读取配置（这里为了方便我写死）

	// fixme 必填: 短信签名，应严格按"签名名称"填写，请参考: https://dysms.console.aliyun.com/dysms.htm#/develop/sign
	signName = "短信签名"

	// fixme 必填: 短信模板Code，应严格按"模板CODE"填写, 请参考: https://dysms.console.aliyun.com/dysms.htm#/develop/template
	verifyTemplateCode = "模板CODE"

	// fixme 必填: 请参阅 https://ak-console.aliyun.com/ 取得您的AK信息
	accessKeyId = "your access key id"
	accessKeySecret = "your access key secret"
}