package mtproto

import (
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

func (s *Client) GenMsgKey(authKey, data []byte) []byte {
	return s.genMsgKey(authKey, data, 0)
}
