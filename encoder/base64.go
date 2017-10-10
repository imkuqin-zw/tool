package encoder

import (
	"bytes"
	"encoding/base64"
	"strings"
)

func Base64Encode(src []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)
	return buf
}

func Base64Decode(src []byte) []byte {
	var encoder *base64.Encoding
	if bytes.Contains(src, []byte{'-'}) {
		encoder = base64.URLEncoding
	} else {
		encoder = base64.StdEncoding
	}

	buf := make([]byte, encoder.DecodedLen(len(src)))
	encoder.Decode(buf, src)
	return buf
}

func Base64EncodeString(src string) (out string) {
	if src == "" {
		return
	}
	byteSrc := []byte(src)

	out = base64.StdEncoding.EncodeToString(byteSrc)

	return
}

func Base64DecodeString(src string) (out string) {
	if src == "" {
		return
	}
	a := len(src) % 4
	if a != 0 {
		src = src + strings.Repeat("=", 4-a)
	}
	var encoder *base64.Encoding
	if strings.Contains(src, "-") {
		encoder = base64.URLEncoding
	} else {
		encoder = base64.StdEncoding
	}
	bs, err := encoder.DecodeString(src)
	if err != nil {
		return
	}

	out = string(bs)
	return
}
