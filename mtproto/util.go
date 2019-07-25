package mtproto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"errors"
	"github.com/imkuqin-zw/tool/file"
	"io"
	"io/ioutil"
	"os"
)

func RandInt128() [16]byte {
	var nonce [16]byte
	_, _ = io.ReadFull(rand.Reader, nonce[:])
	return nonce
}

func RandInt256() [32]byte {
	var nonce [32]byte
	_, _ = io.ReadFull(rand.Reader, nonce[:])
	return nonce
}

func GetAESPubKey(certPath string) (map[uint64]crypto.PublicKey, error) {
	if !file.IsDir(certPath) {
		return nil, errors.New("rsa certificate path is not a directory")
	}
	files, err := file.GetAllFiles(certPath)
	if err != nil {
		return nil, err
	}
	RSAPubKey := make(map[uint64]crypto.PublicKey)
	for i := range files {
		f, err := os.OpenFile(files[i], os.O_RDONLY, 0777)
		if err != nil {
			return nil, err
		}
		pubKey, err := ioutil.ReadAll(f)
		pubKeySha1 := sha1.Sum(pubKey)
		finger := binary.BigEndian.Uint64(pubKeySha1[12:])
		RSAPubKey[finger] = pubKey
		break
	}
	return RSAPubKey, nil
}

func GetAESCert(certPath string) ([]uint64, map[uint64]crypto.PublicKey, map[uint64]crypto.PrivateKey, error) {
	if !file.IsDir(certPath) {
		return nil, nil, nil, errors.New("rsa certificate path is not a directory")
	}
	var (
		fingerArr  = make([]uint64, 0)
		RSAPubKey  = make(map[uint64]crypto.PublicKey)
		RSAPrivKey = make(map[uint64]crypto.PrivateKey)
	)
	files, err := file.GetAllFiles(certPath)
	if err != nil {
		return nil, nil, nil, err
	}
	for i := range files {
		f, err := os.OpenFile(files[i], os.O_RDONLY, 0777)
		if err != nil {
			return nil, nil, nil, err
		}
		var cert []byte
		cert, err = ioutil.ReadAll(f)
		if err != nil {
			return nil, nil, nil, err
		}

		key := bytes.Split(cert, []byte{'#', '#', '#', '#'})
		if len(key) != 2 {
			return nil, nil, nil, errors.New("ase certificate format wrong")
		}
		pubKeySha1 := sha1.Sum(key[0])
		finger := binary.BigEndian.Uint64(pubKeySha1[12:])
		RSAPubKey[finger], RSAPrivKey[finger] = key[0], key[1]
		fingerArr = append(fingerArr, finger)
		break
	}
	return fingerArr, RSAPubKey, RSAPrivKey, nil
}
