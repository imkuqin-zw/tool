package main

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/imkuqin-zw/tool/mtproto"
	"log"
	"net"
	"time"
)

var (
	protoCli *mtproto.Client
	conn     net.Conn
	cliRead  = make(chan []byte, 1)
	cliSend  = make(chan []byte, 10)
	cliVar   = new(CliVar)
)

type CliVar struct {
	cliNonce   []byte
	svrNonce   []byte
	fingers    []uint64
	cliPrivKey crypto.PrivateKey
	cliPubKey  crypto.PublicKey
	svrPubKey  crypto.PublicKey
	newNonce   []byte
	useFinger  uint64
	authKey    []byte
}

func main() {
	var err error
	protoCli, err = mtproto.NewClient("../rsa_cli")
	if err != nil {
		log.Fatal(1, err.Error())
		return
	}
	conn, err = net.Dial("tcp", ":30000")
	if err != nil {
		log.Fatal(2, err.Error())
		return
	}
	go cliReadConn()
	reqPqMulti()
	reqECDH()
	ChangeDH()
}

func cliSendConn(cmd uint32, data []byte) {
	crcId := make([]byte, 4)
	binary.BigEndian.PutUint32(crcId, cmd)
	data = append(crcId, data...)
	dataLen := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLen, uint32(len(data)))
	buf := bytes.NewBuffer(dataLen)
	buf.Write(data)
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}

func cliReadConn() {
	for {
		var lenData [4]byte
		_ = conn.SetReadDeadline(time.Now().Add(time.Second * 20))
		_, err := conn.Read(lenData[:])
		if err != nil {
			_ = conn.Close()
			log.Fatal(3, err.Error())
		}
		dataLen := binary.BigEndian.Uint32(lenData[:])
		data := make([]byte, dataLen)
		_, _ = conn.Read(data)
		cliRead <- data
	}
}

func reqPqMulti() {
	nonce := protoCli.GenNonce128()
	cliSendConn(1, nonce[:])
	rsp := <-cliRead
	data := rsp[4:]
	res := &ResPqMulti{}
	if err := json.Unmarshal(data, res); err != nil {
		log.Fatal(4, err.Error())
		return
	}
	cliVar.svrNonce = res.Nonce
	cliVar.cliNonce = nonce[:]
	cliVar.fingers = res.Finger
	return
}

func reqECDH() {
	var err error
	newNonce := protoCli.GenNewNonce()
	cliVar.newNonce = newNonce[:]
	cliVar.useFinger = protoCli.ChoiceRsaKey(cliVar.fingers)
	params := &ReqECDHParams{
		Finger:   cliVar.useFinger,
		CliNonce: cliVar.cliNonce,
		SvrNonce: cliVar.svrNonce,
	}
	encryptedData := &EncryptedData{
		NewNonce: newNonce[:],
		CliNonce: cliVar.cliNonce,
		SvrNonce: cliVar.svrNonce,
	}
	encrypData, _ := json.Marshal(encryptedData)
	params.EncryptedData, err = protoCli.RsaEncrypt(cliVar.useFinger, protoCli.GetDataWithHash(encrypData))
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	sendData, _ := json.Marshal(params)
	cliSendConn(2, sendData)
	rsp := <-cliRead
	rspData := rsp[4:]
	res := &ResECDH{}
	_ = json.Unmarshal(rspData, res)
	if bytes.Compare(res.SvrNonce, cliVar.svrNonce) != 0 {
		log.Fatal("svr nonce not equal")
		return
	}
	if bytes.Compare(res.CliNonce, cliVar.cliNonce) != 0 {
		log.Fatal("cli nonce not equal")
		return
	}
	data, err := protoCli.AESTmpDecrypt(res.Encrypt, cliVar.newNonce, cliVar.svrNonce)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	data = protoCli.RemoveHash(data)
	encryData := &ResECDHEncrypt{}
	_ = json.Unmarshal(data, encryData)
	if bytes.Compare(encryData.SvrNonce, cliVar.svrNonce) != 0 {
		log.Fatal("svr nonce not equal")
		return
	}
	if bytes.Compare(encryData.CliNonce, cliVar.cliNonce) != 0 {
		log.Fatal("cli nonce not equal")
		return
	}
	cliVar.svrPubKey = &encryData.PubKey
}

func ChangeDH() {
	cliVar.cliPubKey, cliVar.cliPrivKey = protoCli.GenECDHKey()

	req := &ReqChangeDH{
		CliNonce: cliVar.cliNonce,
		SvrNonce: cliVar.svrNonce,
	}
	encry := &ChangeDHEncry{
		CliNonce: cliVar.cliNonce,
		SvrNonce: cliVar.svrNonce,
		PubKey:   *cliVar.cliPubKey.(*[32]byte),
	}
	encryptedData, _ := json.Marshal(encry)
	var err error
	req.EncryptedData, err = protoCli.AESTmpEncrypt(protoCli.GetDataWithHash(encryptedData), cliVar.newNonce, cliVar.svrNonce)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	sendData, _ := json.Marshal(req)
	cliSendConn(3, sendData)
	rsp := <-cliRead
	rspData := rsp[4:]
	res := &ResChangeDH{}
	_ = json.Unmarshal(rspData, res)
	if bytes.Compare(res.SvrNonce, cliVar.svrNonce) != 0 {
		log.Fatal("svr nonce not equal")
		return
	}
	if bytes.Compare(res.CliNonce, cliVar.cliNonce) != 0 {
		log.Fatal("cli nonce not equal")
		return
	}
	cliVar.authKey = protoCli.GenAuthKey(cliVar.cliPrivKey, cliVar.svrPubKey)
	nonceHash := protoCli.GetNewNonceHash(cliVar.newNonce, cliVar.authKey, mtproto.AuthKeyGenOk)
	if !bytes.Contains(res.NonceHash, nonceHash) {
		log.Fatal("new nonce hash not equal")
		return
	}
	fmt.Println("auth_key", hex.EncodeToString(cliVar.authKey))
	//fmt.Println("svrPubKey", string(cliVar.svrPubKey.(*[32]byte)[:]))
	//fmt.Println("cliPubKey", string(cliVar.cliPubKey.(*[32]byte)[:]))
	//fmt.Println("cliPrivKey", string(cliVar.cliPrivKey.(*[32]byte)[:]))
	return
}
