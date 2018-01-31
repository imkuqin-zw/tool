package net_util

import (
	"net/http"
	"time"
	"net/url"
	"strings"
	"io/ioutil"
	"crypto/tls"
)

func Get(url string, params url.Values, header, cookie map[string]string) (body []byte, err error) {
	return request("GET", url, params, header, cookie)
}

func Post(url string, params url.Values, header, cookie map[string]string) (body []byte, err error) {
	return request("POST", url, params, header, cookie)
}

func Put(url string, params url.Values, header, cookie map[string]string) (body []byte, err error) {
	return request("PUT", url, params, header, cookie)
}

func Delete(url string, params url.Values, header, cookie map[string]string) (body []byte, err error) {
	return request("DELETE", url, params, header, cookie)
}

func request(method, url string, params url.Values, header map[string]string, cookie map[string]string) (body []byte, err error) {
	var reader *strings.Reader
	if len(params) > 0 {
		reader = strings.NewReader(params.Encode())
	} else {
		reader = strings.NewReader("")
	}
	req, _ := http.NewRequest(method, url, reader)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	for k, v := range cookie {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}
	// 忽略https证书验证
	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := http.Client{Transport: transport, Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}