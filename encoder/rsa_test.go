package encoder

import (
	"testing"
	"fmt"
)

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDLKa3VECEid2UvpYOjrBtgTxK+vlpJ5S8kVu9NtZ0+/gVMh2Ot
g+xjSYeIVCkaE761MIAx4ZjKe6/Drvvcuux5HPcd7+xq9zCGqHFYOQpIqD245we/
x9HmbOHCIWvyZLpX9ZsmVz0iS9oz4rbMopsvkqsiJuN7hjPGcAOWL5BO4QIDAQAB
AoGACSd6nrQYWh45H/l8Qf66SQ+nD5MyLEw4YJHOPJknWbRGdtlO432jRCIHClyI
cZVcLXve+uBoaw9BrzaOQLbnesT27Z/7Ac/QBWM+RTv7pBBrl88CoOLopE0KDk05
fvpieH1ODve5c13L9JY8hTJXKDlUGHFnzS0UZ8rZuP8kJTUCQQD2mpG5kRfV150f
qsqk44nqwHZ3hNyNP5tGrtD3IfSB5SQqn4a5uHesk0kD02ns3gHPLc5DCyr4WUeD
FoK4dCgLAkEA0udgWrZeanFXhduT4p4zroV5sUWxP3eTC0svpTwwspfsHlUm9aWw
JiQe9KN01yvVusrvt8jVSgKlIl696Z38QwJBAOtycKfn7AXzssTFYG1GAivsTi+W
3qzNigdWaZVLChPrHzjCzvMLONfAV/obJAtPfBK+/Svtwb0UIL78AxrxbDkCQGyM
oAAwawn4CicgK85wxILnugmuqBrVbX5blUtDXoEdRm8aPrIiNDZ9Ut1xH9r7ecbp
WlZKbNTp5Zp6Dt8l7EcCQQDigS5g9RbRZPSJaCHKq33SpYV1B961IFLAWWLXBY+y
2p+35NlUgDPqtGcQeQbABT9Wo4L7UOBAp8/6hW3RN3dH
-----END RSA PRIVATE KEY-----
`)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLKa3VECEid2UvpYOjrBtgTxK+
vlpJ5S8kVu9NtZ0+/gVMh2Otg+xjSYeIVCkaE761MIAx4ZjKe6/Drvvcuux5HPcd
7+xq9zCGqHFYOQpIqD245we/x9HmbOHCIWvyZLpX9ZsmVz0iS9oz4rbMopsvkqsi
JuN7hjPGcAOWL5BO4QIDAQAB
-----END PUBLIC KEY-----
`)

func TestRsaEncryptAndDecrypt(t *testing.T) {
	out, err := RsaEncrypt([]byte("test123456"), publicKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(out))
	out, err = RsaDecrypt(out, privateKey)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(out))
}

func TestGeneratePrivKeyAndPubKey(t *testing.T) {
	private, publicKey, err := GeneratePrivKeyAndPubKey(1024)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(private))
	fmt.Println(string(publicKey))

}