package mtproto

type Client struct {
	common
}

func (c *Client) CreateNewNonce() [32]byte {
	return RandInt256()
}




