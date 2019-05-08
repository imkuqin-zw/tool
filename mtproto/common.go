package mtproto

import (
	"crypto"
	"crypto/rand"
	"crypto/sha1"
	"github.com/imkuqin-zw/tool/encoder"
)

type common struct {
	Ecdh      encoder.ECDH
	RSAPubKey map[uint64]crypto.PublicKey
}

func (c *common) AESTmpEncrypt(data, newNonce, serverNonce []byte) ([]byte, error) {
	key, iv := DeriveTmpAESKeyIV(newNonce, serverNonce)
	return encoder.AESCBCEncrypt(data, key, iv)
}

func (c *common) AESTmpDecrypt(data, newNonce, serverNonce []byte) ([]byte, error) {
	key, iv := DeriveTmpAESKeyIV(newNonce, serverNonce)
	return encoder.AESCBCDecrypt(data, key, iv)
}

func (c *common) CreateECDHKey() (crypto.PublicKey, crypto.PublicKey) {
	pubKey, prvKey, _ := c.Ecdh.GenerateKey(rand.Reader)
	return pubKey, prvKey
}

func (c *common) RsaEncrypt(finger uint64, data []byte) ([]byte, error) {
	return encoder.RsaEncrypt(data, c.RSAPubKey[finger].([]byte))
}

func (c *common) CreateNonce128() [16]byte {
	return RandInt128()
}

func (c *common) GenerateSharedKey(privKey crypto.PrivateKey, pubKey crypto.PublicKey) []byte {
	secret, _ := c.Ecdh.GenerateSharedSecret(privKey, pubKey)
	return secret
}

func (c *common) GetAuthKeyId(authKey []byte) []byte {
	temp := sha1.Sum(authKey)
	return temp[12:]
}
