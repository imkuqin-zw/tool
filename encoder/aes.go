package encoder

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func AESCBCDecrypt(data, key, iv []byte) ([]byte, error) {
	ciph, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(data)%ciph.BlockSize() != 0 {
		return nil, errors.New("input not multiple of block size")
	}

	bs := ciph.BlockSize()
	if len(iv) < bs {
		return nil, errors.New("IV length must equal block size")
	}
	dst := make([]byte, len(data))
	iv16bytes := iv[:bs]
	mode := cipher.NewCBCDecrypter(ciph, iv16bytes)
	mode.CryptBlocks(dst, data)
	dst, err = PKCS5UnPadding(dst)
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func AESCBCEncrypt(data, key, iv []byte) ([]byte, error) {
	ciph, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := ciph.BlockSize()
	if len(iv) < bs {
		return nil, errors.New("IV length must equal block size")
	}
	iv16bytes := iv[:bs]
	data = PKCS5Padding(data, bs)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(data))
	blockModel := cipher.NewCBCEncrypter(ciph, iv16bytes)
	blockModel.CryptBlocks(dst, data)
	return dst, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	paddingLen := int(origData[length-1])
	if paddingLen >= length {
		return nil, errors.New("padding size error")
	}
	return origData[:length-paddingLen], nil
}
