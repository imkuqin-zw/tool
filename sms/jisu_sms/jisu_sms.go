package jisu_sms

import (
	"github.com/imkuqin-zw/tool/order"
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
	"encoding/json"
)

var (
	domain = ""
	key = ""
	signName = ""
	verifyTemplate = ""
)

type SendResponse struct {
	Status 	string
	Msg		string
	Result	*RespResult
}

type RespResult struct {
	Accountid 	string
	Msgid 		string
	Count		string
}

const EXPIRE = 1800

func SendVerifySms(phone string, code string) (err error) {
	kv := order.NewKVMap(3)
	kv.Set("appkey", key)
	kv.Set("mobile", phone)
	kv.Set("content", fmt.Sprintf(verifyTemplate, code, signName))
	params := url.Values{
		"appkey": []string{key},
		"mobile": []string{phone},
		"content": []string{fmt.Sprintf(verifyTemplate, code, signName)},
	}

	resp, err := http.Get(domain + "/?" + params.Encode())
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result SendResponse
	if err = json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("[短信发送失败]%s", "解析返回数据错误")
	}
	if result.Status != "0" {
		return fmt.Errorf("[短信发送失败]%s", result.Msg)
	}
	return
}


func init() {
	domain = ""
	key = ""
	signName = ""
	verifyTemplate = "您的验证码是%s,请勿转发或告知他人,30分钟内有效。请在页面中输入以完成验证【%s】"
}