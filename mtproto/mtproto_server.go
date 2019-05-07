package mtproto

import "crypto"

type Server struct {
	ASEList map[string]crypto.PrivateKey
}

func (s *Server) CreateSvrNonce() [16]byte {
	return RandInt128()
}
