package mtproto

import (
	"bytes"
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

//auth_key为第1024直接
func (c *Client) GetMsgKey(authKey, data []byte) []byte {
	buf := bytes.NewBuffer(authKey[88:120])
	buf.Write(data)
	tmp := sha256.Sum256(buf.Bytes())
	return tmp[8:24]
}

func (c *Client) GetAesKeyIv(auth_key, msgKey) {

}