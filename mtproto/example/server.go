package main

import (
	"bytes"
	"context"
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
	svrRead     = make(chan []byte, 1)
	svrSend     = make(chan []byte, 10)
	ctx, cancel = context.WithCancel(context.Background())
	protoSvr    *mtproto.Server
	svrVar      = new(SvrVar)
)

type SvrVar struct {
	cliNonce   []byte
	svrNonce   []byte
	newNonce   []byte
	cliPubKey  crypto.PublicKey
	svrPubKey  crypto.PublicKey
	svrPrivKey crypto.PrivateKey
	authKey    []byte
}

func main() {
	var err error
	protoSvr, err = mtproto.NewServer("../rsa_cert")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":30000")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	t, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	go func() {
		for {
			conn, err := t.AcceptTCP()
			if err != nil || conn == nil {
				log.Println(err.Error())
				continue
			}
			go server(conn)
		}
	}()
	select {}
}

func server(conn *net.TCPConn) {
	go readConn(conn)
	go sendConn(conn)
	go handel()
}

func readConn(conn *net.TCPConn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			var lenData [4]byte
			_ = conn.SetReadDeadline(time.Now().Add(time.Second * 20))
			_, err := conn.Read(lenData[:])
			if err != nil {
				svrRead = make(chan []byte, 1)
				return
			}
			dataLen := binary.BigEndian.Uint32(lenData[:])
			data := make([]byte, dataLen)
			_, _ = conn.Read(data)
			svrRead <- data
		}

	}
}

func sendConn(conn *net.TCPConn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			data := <-svrSend
			dataLen := len(data)
			start := 0
			for start < dataLen {
				n, err := conn.Write(data[start:])
				if err != nil {
					svrSend = make(chan []byte, 10)
					return
				}
				start += n
			}
		}
	}
}

func sendData(data []byte) {
	dataLen := make([]byte, 4)
	binary.BigEndian.PutUint32(dataLen, uint32(len(data)))
	buf := bytes.NewBuffer(dataLen)
	buf.Write(data)
	svrSend <- buf.Bytes()
}

func handel() {
	for {
		select {
		case data := <-svrRead:
			cmd := binary.BigEndian.Uint32(data[:4])
			switch cmd {
			case 1:
				reqPqMultiHandle(data)
			case 2:
				reqECDHHandle(data)
			case 3:
				ChangeDHHandle(data)
			default:
				log.Println("cmd not found", cmd)
			}
		case <-ctx.Done():
			return
		}

	}
}

func reqPqMultiHandle(req []byte) {
	svrVar.cliNonce = req[4:]
	nonce := protoSvr.GenNonce128()
	svrVar.svrNonce = nonce[:]
	res := ResPqMulti{
		Nonce:  nonce[:],
		Finger: protoSvr.GetRsaKeyFingers(),
	}
	data, _ := json.Marshal(res)
	sendData(append(req[:4], data...))
}

func reqECDHHandle(req []byte) {
	params := &ReqECDHParams{}
	_ = json.Unmarshal(req[4:], params)
	if 0 != bytes.Compare(params.SvrNonce, svrVar.svrNonce) {
		log.Fatal("svr nonce not equal")
		return
	}
	if 0 != bytes.Compare(params.CliNonce, svrVar.cliNonce) {
		log.Fatal("cli nonce not equal")
		return
	}
	encryData, _ := protoSvr.RsaDecrypt(params.Finger, params.EncryptedData)
	if !protoSvr.CheckHash(encryData) {
		log.Fatal("hash not equal")
		return
	}
	encryptedData := &EncryptedData{}
	_ = json.Unmarshal(protoSvr.RemoveHash(encryData), encryptedData)
	if 0 != bytes.Compare(encryptedData.SvrNonce, svrVar.svrNonce) {
		log.Fatal("svr nonce not equal")
		return
	}
	if 0 != bytes.Compare(encryptedData.CliNonce, svrVar.cliNonce) {
		log.Fatal("cli nonce not equal")
		return
	}
	svrVar.newNonce = encryptedData.NewNonce
	svrVar.svrPubKey, svrVar.svrPrivKey = protoSvr.GenECDHKey()
	//svrVar.cliPubKey = encryptedData.pubKey
	//authKey := protoSvr.GenSharedKey(svrVar.svrPrivKey, svrVar.cliPubKey)
	res := &ResECDHEncrypt{
		CliNonce: svrVar.cliNonce,
		SvrNonce: svrVar.svrNonce,
		PubKey:   *svrVar.svrPubKey.(*[32]byte),
	}
	data, _ := json.Marshal(res)
	hashData := protoSvr.GetDataWithHash(data)
	spt := &ResECDH{
		CliNonce: svrVar.cliNonce,
		SvrNonce: svrVar.svrNonce,
	}
	var err error
	spt.Encrypt, err = protoSvr.AESTmpEncrypt(hashData, svrVar.newNonce, svrVar.svrNonce)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	data, _ = json.Marshal(spt)
	sendData(append(req[:4], data...))
}

func ChangeDHHandle(req []byte) {
	params := &ReqChangeDH{}
	_ = json.Unmarshal(req[4:], params)
	if 0 != bytes.Compare(params.SvrNonce, svrVar.svrNonce) {
		log.Fatal("55 svr nonce not equal")
		return
	}
	if 0 != bytes.Compare(params.CliNonce, svrVar.cliNonce) {
		log.Fatal("cli nonce not equal")
		return
	}
	encryData, _ := protoSvr.AESTmpDecrypt(params.EncryptedData, svrVar.newNonce, svrVar.svrNonce)
	if !protoSvr.CheckHash(encryData) {
		log.Fatal("hash not equal")
		return
	}
	encryptedData := &ChangeDHEncry{}
	_ = json.Unmarshal(protoSvr.RemoveHash(encryData), encryptedData)
	if 0 != bytes.Compare(encryptedData.SvrNonce, svrVar.svrNonce) {
		log.Fatal("66 svr nonce not equal")
		return
	}
	if 0 != bytes.Compare(encryptedData.CliNonce, svrVar.cliNonce) {
		log.Fatal("cli nonce not equal")
		return
	}
	svrVar.cliPubKey = &encryptedData.PubKey
	svrVar.authKey = protoSvr.GenAuthKey(svrVar.svrPrivKey, svrVar.cliPubKey)
	spt := &ResChangeDH{
		CliNonce:  svrVar.cliNonce,
		SvrNonce:  svrVar.svrNonce,
		NonceHash: protoSvr.GetNewNonceHash(svrVar.newNonce, svrVar.authKey, mtproto.AuthKeyGenOk),
	}
	data, _ := json.Marshal(spt)
	sendData(append(req[:4], data...))
	fmt.Println("auth_key", hex.EncodeToString(svrVar.authKey))
	//fmt.Println("svrPubKey", string(svrVar.svrPubKey.(*[32]byte)[:]))
	//fmt.Println("vliPubKey", string(svrVar.cliPubKey.(*[32]byte)[:]))
	//fmt.Println("svrPrivKey", string(svrVar.svrPrivKey.(*[32]byte)[:]))
}
