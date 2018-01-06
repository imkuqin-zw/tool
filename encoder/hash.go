package encoder

import (
	"crypto/sha256"
	"encoding/hex"
	"crypto/sha1"
)

func Sha256(src string) (result string, err error) {
	hash := sha256.New()
	_, err = hash.Write([]byte(src))
	if err != nil {
		return
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return
}

func Sha1(src string) (result string, err error) {
	hash := sha1.New()
	_, err = hash.Write([]byte(src))
	if err != nil {
		return
	}
	result = hex.EncodeToString(hash.Sum(nil))
	return
}