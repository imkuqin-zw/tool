package encoder

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Byte(src []byte) (result []byte, err error) {
	hash := md5.New()
	_, err = hash.Write([]byte(src))
	if err != nil {
		return
	}
	md5Byte := hash.Sum(nil)
	result = make([]byte, hex.EncodedLen(len(md5Byte)))
	hex.Encode(result, md5Byte)

	return
}

func Md5String(src string) (result string, err error) {
	hash := md5.New()
	_, err = hash.Write([]byte(src))
	if err != nil {
		return
	}
	md5Byte := hash.Sum(nil)
	result = hex.EncodeToString(md5Byte)
	return
}
