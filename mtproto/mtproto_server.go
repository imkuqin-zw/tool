package mtproto

import (
	"crypto"
	"github.com/imkuqin-zw/tool/encoder"
)

type Server struct {
	common
	RSAPrivKey map[uint64]crypto.PrivateKey
}

func NewServer(certPath string) (*Server, error) {
	var svr = &Server{}
	svr.Ecdh = encoder.NewCurve25519ECDH()
	var err error
	svr.RSAPubKey, svr.RSAPrivKey, err = GetAESCert(certPath)
	if err != nil {
		return nil, err
	}
	return svr, nil
}

func (s *Server) CreateSvrNonce() [16]byte {
	return RandInt128()
}

func (s *Server) RsaDecrypt(finger uint64, data []byte) ([]byte, error) {
	return encoder.RsaDecrypt(data, s.RSAPrivKey[finger].([]byte))
}

func (s *Server) GenMsgKey(authKey, data []byte) []byte {
	return s.genMsgKey(authKey, data, 8)
}
