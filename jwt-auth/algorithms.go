package jwt_auth

import "tool/encoder"

type Algorithms struct{}

func (c *Algorithms) Sha256(headerStr string, payloadStr string, secret string) (result string, err error) {
	result, err = encoder.Sha256(headerStr + payloadStr + secret)
	return
}

func (c *Algorithms) Md5(headerStr string, payloadStr string, secret string) (result string, err error) {
	result, err = encoder.Md5String(headerStr + payloadStr + secret)
	return
}
