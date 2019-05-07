package mtproto

import (
	"crypto"
	"github.com/imkuqin-zw/tool/encoder"
)

type Server struct {
	common
	RSAPrivKey map[uint64]crypto.PrivateKey
}

func (s *Server) CreateSvrNonce() [16]byte {
	return RandInt128()
}

func (s *Server) RsaDecrypt(finger uint64, data []byte) ([]byte, error) {
	return encoder.RsaDecrypt(data, s.RSAPrivKey[finger].([]byte))
}