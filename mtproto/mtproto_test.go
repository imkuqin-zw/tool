package mtproto

import (
	"fmt"
	"testing"
)

func Test_MTProto(t *testing.T) {
	data, _ := GetAESPubKey("./rsa_cli")
	for i, item := range data {
		fmt.Println(i)
		fmt.Println(string(item.([]byte)))
	}
	fmt.Println("-----------------------------------------")
	_, pubkey, _, _ := GetAESCert("./rsa_cert")
	for i, item := range pubkey {
		fmt.Println(i)
		fmt.Println(string(item.([]byte)))
	}
}
