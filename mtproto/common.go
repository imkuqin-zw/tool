package mtproto

import (
	"crypto"
	"crypto/rand"
	"errors"
	"github.com/imkuqin-zw/tool/encoder"
	"github.com/imkuqin-zw/tool/file"
)

type common struct {
	ecdh      encoder.ECDH
	RSAPubKey map[uint64]crypto.PublicKey
}

func newCommon(certPath string) (comm common, err error) {
	comm = common{
		ecdh:      encoder.NewCurve25519ECDH(),
		RSAPubKey: make(map[uint64]crypto.PublicKey),
	}
	if !file.IsDir(certPath) {
		err = errors.New("rsa certificate path is not a directory")
		return
	}

}

func (c *common) AESEncrypt(data, newNonce, serverNonce []byte) ([]byte, error) {
	key, iv := DeriveTmpAESKeyIV(newNonce, serverNonce)
	return encoder.AESCBCEncrypt(data, key, iv)
}

func (c *common) AESDecrypt(data, newNonce, serverNonce []byte) ([]byte, error) {
	key, iv := DeriveTmpAESKeyIV(newNonce, serverNonce)
	return encoder.AESCBCDecrypt(data, key, iv)
}

func (c *common) CreateECDHKey() (crypto.PublicKey, crypto.PublicKey) {
	pubKey, prvKey, _ := c.ecdh.GenerateKey(rand.Reader)
	return pubKey, prvKey
}

func (c *common) RsaEncrypt(finger uint64, data []byte) ([]byte, error) {
	return encoder.RsaEncrypt(data, c.RSAPubKey[finger].([]byte))
}

func (c *common) CreateNonce128() [16]byte {
	return RandInt128()
}

func (c *common) GenerateSharedKey(privKey crypto.PrivateKey, pubKey crypto.PublicKey) []byte {
	secret, _ := c.ecdh.GenerateSharedSecret(privKey, pubKey)
	return secret
}
