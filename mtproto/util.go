package mtproto

import (
	"crypto/rand"
	"io"
)

func RandInt128() [16]byte {
	var nonce [16]byte
	_, _ = io.ReadFull(rand.Reader, nonce[:])
	return nonce
}

func RandInt256() [32]byte {
	var nonce [32]byte
	_, _ = io.ReadFull(rand.Reader, nonce[:])
	return nonce
}
