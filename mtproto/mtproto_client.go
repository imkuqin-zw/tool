package mtproto

import "github.com/imkuqin-zw/tool/encoder"

type Client struct {
	common
}

func NewClient() *Client {
	return &Client{
		common{
			ecdh:encoder.NewCurve25519ECDH(),
		},
	}
}

func (c *Client) CreateNewNonce() [32]byte {
	return RandInt256()
}




