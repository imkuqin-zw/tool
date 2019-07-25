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

func (s *Client) GenMsgKey(authKey, data []byte) []byte {
	return s.genMsgKey(authKey, data, 0)
}

func (s *Client) ChoiceRsaKey(fingers []uint64) uint64 {
	return fingers[rand.Intn(len(fingers))]
}
