package net_util

import (
	"testing"
	"fmt"
	"net/url"
)

func TestGet(t *testing.T)  {
	body, err := Get("http://www.baidu.com", url.Values{}, nil, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(body))
	return
}
