package jwt_auth

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"
	"tool/encoder"
)

var (
	ErrEncHeader = errors.New("[JWT] encode header error")
	ErrDecHeader = errors.New("[JWT] decode header error")
	ErrEncPayload = errors.New("[JWT] encode payload error")
	ErrDecPayload = errors.New("[JWT] decode payload error")
	ErrEncSign = errors.New("[JWT] encode sign error")
	ErrFormat = errors.New("[JWT] token format error")
	ErrTokenSign = errors.New("[JWT] token signature error")
	ErrTokenExpired = errors.New("[JWT] token signature error")
)

func encodeHeader(alg string) (headerStr string, err error) {
	header := Header{Typ: "JWT", Alg: alg}
	kv, err := json.Marshal(header)
	if err != nil {
		return
	}

	headerStr = encoder.Base64EncodeString(string(kv))
	if headerStr == "" {
		err = ErrEncHeader
		return
	}
	return
}

func decodeHeader(src string) (header Header, err error) {
	headerStr := encoder.Base64DecodeString(src)
	header = Header{}
	err = json.Unmarshal([]byte(headerStr), &header)
	if err != nil {
		err = ErrDecHeader
		return
	}

	return
}

func encodePayload(m *Payload) (payloadStr string, err error) {
	m.Iat = time.Now().Unix()
	kv, err := json.Marshal(m)
	if err != nil {
		return
	}

	payloadStr = encoder.Base64EncodeString(string(kv))
	if payloadStr == "" {
		err = ErrEncPayload
		return
	}

	return
}

func decodePayload(src string) (payload Payload, err error) {
	payloadStr := encoder.Base64DecodeString(src)
	payload = Payload{}
	err = json.Unmarshal([]byte(payloadStr), &payload)
	if err != nil {
		err = ErrDecPayload
		return
	}

	return
}

func encodeSign(alg string, secret string, headerStr string, payloadStr string) (signStr string, err error) {
	algorithms := &Algorithms{}
	v := reflect.ValueOf(algorithms)
	params := make([]reflect.Value, 3)
	params[0] = reflect.ValueOf(headerStr)
	params[1] = reflect.ValueOf(secret)
	params[2] = reflect.ValueOf(payloadStr)
	response := v.MethodByName(alg).Call(params)

	if response[1].Interface() != nil {
		err = ErrEncSign
		return
	}
	signStr = response[0].Interface().(string)
	return
}

func EncodeJWT(alg string, iss string, sub string, aud string, exp int64, nbf int64, secret string) (out string, err error) {
	headerStr, err := encodeHeader(alg)
	if err != nil {
		return
	}

	payload := &Payload{
		Iss: iss,
		Sub: sub,
		Aud: aud,
		Exp: exp,
		Nbf: nbf,
	}
	payloadStr, err := encodePayload(payload)
	if err != nil {
		return
	}
	signStr, err := encodeSign(alg, secret, headerStr, payloadStr)
	if err != nil {
		return
	}
	out = headerStr + "." + payloadStr + "." + signStr

	return
}

func DecodeJwT(src string, secret string) (result string, err error) {
	jwtStrs := strings.Split(src, ".")
	if len(jwtStrs) != 3 {
		err = ErrFormat
		return
	}

	header, err := decodeHeader(jwtStrs[0])
	if err != nil {
		return
	}

	signStr, err := encodeSign(header.Alg, secret, jwtStrs[0], jwtStrs[1])
	if err != nil {
		return
	}

	if signStr != jwtStrs[2] {
		err = ErrTokenSign
		return
	}

	payload, err := decodePayload(jwtStrs[1])
	if err != nil {
		return
	}

	if payload.Exp < time.Now().Unix() {
		err = ErrTokenExpired
		return
	}

	result = payload.Sub
	return
}
