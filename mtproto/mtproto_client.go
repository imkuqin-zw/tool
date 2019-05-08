package mtproto

import (
	"crypto/sha256"
	"github.com/imkuqin-zw/tool/encoder"
)

type Client struct {
	common
}

func NewClient(certPath string) (*Client, error) {
	var cli = &Client{}
	cli.Ecdh = encoder.NewCurve25519ECDH()
	var err error
	cli.RSAPubKey, err = GetAESPubKey(certPath)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func (c *Client) CreateNewNonce() [32]byte {
	return RandInt256()
}

func (c *Client) GetMsgKey(authKey, data []byte) []byte {
	authKey[88]
	msgKey := sha256.Sum256()
}
