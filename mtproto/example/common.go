package main

type ReqChangeDH struct {
	CliNonce      []byte
	SvrNonce      []byte
	EncryptedData []byte
}

type ChangeDHEncry struct {
	CliNonce []byte
	SvrNonce []byte
	RetryId  uint64
	PubKey   [32]byte
}

type ReqECDHParams struct {
	Finger        uint64
	CliNonce      []byte
	SvrNonce      []byte
	EncryptedData []byte
}

type EncryptedData struct {
	NewNonce []byte
	CliNonce []byte
	SvrNonce []byte
}

type ResECDHEncrypt struct {
	CliNonce []byte
	SvrNonce []byte
	PubKey   [32]byte
}

type ResECDH struct {
	CliNonce []byte
	SvrNonce []byte
	Encrypt  []byte
}

type ResPqMulti struct {
	Nonce  []byte
	Finger []uint64
}

type ResChangeDH struct {
	CliNonce  []byte
	SvrNonce  []byte
	NonceHash []byte
}
