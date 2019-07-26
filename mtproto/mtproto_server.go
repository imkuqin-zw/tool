package mtproto

import (
	"crypto"
	"github.com/imkuqin-zw/tool/encoder"
)

type Server struct {
	common
	fingers    []uint64
	RSAPrivKey map[uint64]crypto.PrivateKey
}

func NewServer(certPath string) (*Server, error) {
	var svr = &Server{}
	svr.Ecdh = encoder.NewCurve25519ECDH()
	//svr.Ecdh = encoder.NewEllipticECDH(elliptic.P521())
	var err error
	svr.fingers, svr.RSAPubKey, svr.RSAPrivKey, err = GetAESCert(certPath)
	if err != nil {
		return nil, err
	}
	return svr, nil
}

func (s *Server) GetRsaKeyFingers() []uint64 {
	return s.fingers
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

func (s *Server) GetInitSalt(newNonce, svrNonce []byte) []byte {
	return append(newNonce[:8], svrNonce[:8]...)
}

func (s *Server) AESEncrypt(authKey, msgKey, data []byte) ([]byte, error) {
	key, iv := s.genAesKeyIv(authKey, msgKey, 8)
	return encoder.AESCBCEncrypt(data, key, iv)
}

func (s *Server) AESDecrypt(authKey, msgKey, data []byte) ([]byte, error) {
	key, iv := s.genAesKeyIv(authKey, msgKey, 0)
	return encoder.AESCBCDecrypt(data, key, iv)
}
