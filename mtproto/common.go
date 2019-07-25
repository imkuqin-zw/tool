package mtproto

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"github.com/imkuqin-zw/tool/encoder"
)

type common struct {
	Ecdh      encoder.ECDH
	RSAPubKey map[uint64]crypto.PublicKey
}

func (c *common) deriveTmpAESKeyIV(newNonce []byte, serverNonce []byte) ([]byte, []byte) {
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

func (c *common) AESTmpEncrypt(data, newNonce, serverNonce []byte) ([]byte, error) {
	key, iv := c.deriveTmpAESKeyIV(newNonce, serverNonce)
	return encoder.AESCBCEncrypt(data, key, iv)
}

func (c *common) AESTmpDecrypt(data, newNonce, serverNonce []byte) ([]byte, error) {
	key, iv := c.deriveTmpAESKeyIV(newNonce, serverNonce)
	return encoder.AESCBCDecrypt(data, key, iv)
}

func (c *common) GenECDHKey() (crypto.PublicKey, crypto.PrivateKey) {
	prvKey, pubKey, _ := c.Ecdh.GenerateKey(rand.Reader)
	return pubKey, prvKey
}

func (c *common) RsaEncrypt(finger uint64, data []byte) ([]byte, error) {
	return encoder.RsaEncrypt(data, c.RSAPubKey[finger].([]byte))
}

func (c *common) GenNonce128() [16]byte {
	return RandInt128()
}

func (c *common) GenAuthKey(privKey crypto.PrivateKey, pubKey crypto.PublicKey) []byte {
	secret, _ := c.Ecdh.GenerateSharedSecret(privKey, pubKey)
	return secret
}

func (c *common) GetAuthKeyId(authKey []byte) []byte {
	authKeySha1 := sha1.Sum(authKey)
	return authKeySha1[12:]
}

//其中x = 0表示从客户端到服务器的消息，x = 8表示从服务器到客户端的消息
func (c *common) genAesKeyIv(authKey, msgKey []byte, x int) (key []byte, iv []byte) {
	buf := bytes.NewBuffer(msgKey)
	buf.Write(authKey[x : 36+x])
	a := sha256.Sum256(buf.Bytes())

	buf.Reset()
	buf.Write(authKey[40+x : 76+x])
	buf.Write(msgKey)
	b := sha256.Sum256(buf.Bytes())

	buf.Reset()
	buf.Write(a[:8])
	buf.Write(b[8:24])
	buf.Write(a[24:])
	key = buf.Bytes()

	buf.Reset()
	buf.Write(b[:8])
	buf.Write(a[8:24])
	buf.Write(b[24:])
	iv = buf.Bytes()
	return
}

func (c *common) GenCliAesKeyIv(authKey, msgKey []byte) ([]byte, []byte) {
	return c.genAesKeyIv(authKey, msgKey, 0)
}

func (c *common) GenSvrAesKeyIv(authKey, msgKey []byte) ([]byte, []byte) {
	return c.genAesKeyIv(authKey, msgKey, 8)
}

func (c *common) genMsgKey(authKey, data []byte, x int) []byte {
	buf := bytes.NewBuffer(authKey[88+x : 120+x])
	buf.Write(data)
	tmp := sha256.Sum256(buf.Bytes())
	return tmp[8:24]
}

//func (c *common) Sha1(data []byte) []byte {
//	hash := sha1.New()
//	hash.Write(data)
//	return hash.Sum(nil)
//}

// SHA1（数据）+ 数据
func (c *common) GetDataWithHash(data []byte) []byte {
	hash := sha1.Sum(data)
	return append(hash[:], data...)
}

func (c *common) CheckHash(data []byte) bool {
	hash := sha1.Sum(data[20:])
	return 0 == bytes.Compare(hash[:], data[:20])
}

func (c *common) RemoveHash(data []byte) []byte {
	return data[20:]
}

func (s *common) GetAuthKeyAuxHash(authKey []byte) []byte {
	hash := sha1.Sum(authKey)
	return hash[:8]
}

func (s *common) GetNewNonceHash1(newNonce, authKey []byte) []byte {
	buf := bytes.NewBuffer(newNonce)
	buf.Write([]byte{1})
	buf.Write(s.GetAuthKeyAuxHash(authKey))
	hash := sha1.Sum(buf.Bytes())
	return hash[4:]
}
