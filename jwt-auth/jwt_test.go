package jwt_auth

import (
	"fmt"
	"testing"
	"time"
)

func Test_EncodeJWT(t *testing.T) {
	token, err := EncodeJWT("Sha256", "zw", "5", "zw", time.Now().Unix()+3600, time.Now().Unix(), "asdfghjk12314rewtrecvzcvgfg")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(token)
}

func Test_DecodeJwt(t *testing.T) {
	str, err := DecodeJwT("eyJ0eXAiOiJKV1QiLCJhbGciOiJTaGEyNTYifQ==.eyJpc3MiOiJ6dyIsInN1YiI6IjUiLCJhdWQiOiJ6dyIsImV4cCI6MTUwMTQwODMyOSwibmJmIjoxNTAxNDA0NzI5LCJpYXQiOjE1MDE0MDQ3Mjl9.ab70acc2c15a4be211def076fab1138ede2dbed5af9400bc21afab37582396cd", "asdfghjk12314rewtrecvzcvgfg")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(str)
}
