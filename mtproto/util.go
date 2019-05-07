package mtproto

import (
	"crypto/rand"
	"crypto/sha1"
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

func DeriveTmpAESKeyIV(newNonce []byte, serverNonce []byte) ([]byte, []byte) {
	key := make([]byte, 32)
	iv := make([]byte, 32)

	src := make([]byte, 64)
	copy(src[:32], newNonce)
	copy(src[32:48], serverNonce)
	ns := sha1.Sum(src[:48])

	copy(src[:32], newNonce)
	copy(src[32:64], newNonce)
	nn := sha1.Sum(src[:64]) //

	//tmp_aes_key := SHA1(new_nonce + server_nonce) + substr (SHA1(server_nonce + new_nonce), 0, 12);
	copy(key[:20], ns[:])
	copy(key[20:], ns[:12])

	// tmp_aes_iv := substr (SHA1(server_nonce + new_nonce), 12, 8) + SHA1(new_nonce + new_nonce) + substr (new_nonce, 0, 4);
	copy(iv[:8], ns[12:20])
	copy(iv[8:28], nn[:])
	copy(iv[28:], newNonce[:4])

	return key, iv
}
