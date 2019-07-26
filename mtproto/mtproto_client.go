package mtproto

import (
	"github.com/imkuqin-zw/tool/encoder"
	"math/rand"
)

type Client struct {
	common
}

func NewClient(certPath string) (*Client, error) {
	var cli = &Client{}
	cli.Ecdh = encoder.NewCurve25519ECDH()
	//cli.Ecdh = encoder.NewEllipticECDH(elliptic.P521())
	var err error
	cli.RSAPubKey, err = GetAESPubKey(certPath)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (c *Client) GenNewNonce() [32]byte {
	return RandInt256()
}

func (c *Client) GenMsgKey(authKey, data []byte) []byte {
	return c.genMsgKey(authKey, data, 0)
}

func (c *Client) ChoiceRsaKey(fingers []uint64) uint64 {
	return fingers[rand.Intn(len(fingers))]
}

func (c *Client) AESEncrypt(authKey, msgKey, data []byte) ([]byte, error) {
	key, iv := c.genAesKeyIv(authKey, msgKey, 0)
	return encoder.AESCBCEncrypt(data, key, iv)
}

func (c *Client) AESDecrypt(authKey, msgKey, data []byte) ([]byte, error) {
	key, iv := c.genAesKeyIv(authKey, msgKey, 8)
	return encoder.AESCBCDecrypt(data, key, iv)
}
