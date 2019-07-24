package main

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"encoding/json"
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
	pubKey   []byte
}

func main() {
	var err error
	protoCli, err = mtproto.NewClient("")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	conn, err = net.Dial("tcp", "0.0.0.0:30000")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	go cliReadConn()
	reqPqMulti()
}

func cliSendConn(cmd uint32, data []byte) {
	crcId := make([]byte, 0, 4)
	binary.BigEndian.PutUint32(crcId, cmd)
	data = append(crcId, data...)
	dataLen := make([]byte, 0, 4)
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
			log.Fatal(err.Error())
		}
		dataLen := binary.BigEndian.Uint32(lenData[:])
		data := make([]byte, 0, dataLen)
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
		log.Fatal(err.Error())
		return
	}
	cliVar.svrNonce = res.Nonce
	cliVar.cliNonce = nonce[:]
	cliVar.fingers = res.Finger
	return
}

func reqECDH() {
	var err error
	cliVar.cliPubKey, cliVar.cliPrivKey = protoCli.GenECDHKey()
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
		pubKey:   cliVar.cliPubKey.([]byte),
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
	encryData := &ResECDHEncrypt{}
	_ = json.Unmarshal(data, encryData)
	if bytes.Compare(res.SvrNonce, cliVar.svrNonce) != 0 {
		log.Fatal("svr nonce not equal")
		return
	}
	if bytes.Compare(res.CliNonce, cliVar.cliNonce) != 0 {
		log.Fatal("cli nonce not equal")
		return
	}
	cliVar.svrPubKey = encryData.Pubkey
}
