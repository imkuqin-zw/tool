package mtproto

import (
	"crypto"
	"crypto/rand"
	"github.com/imkuqin-zw/tool/encoder"
)

type Client struct {
	ASEList map[uint64]crypto.PublicKey
	ecdh    encoder.ECDH
}

func (c *Client) CreateCliNonce() [16]byte {
	return RandInt128()
}

func (c *Client) CreateNewNonce() [32]byte {
	return RandInt256()
}

func (c *Client) CreateECDHKey() (crypto.PublicKey, crypto.PublicKey) {
	pubKey, prvKey, _ := c.ecdh.GenerateKey(rand.Reader)
	return pubKey, prvKey
}

func (c *Client) RsaEncrypt(finger uint64, data []byte) ([]byte, error) {
	return encoder.RsaEncrypt(data, c.ASEList[finger].([]byte))
}
