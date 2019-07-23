package main

import (
	"bytes"
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
	cliNonce []byte
	fingers  []uint64
}

type reqECDHParams struct {
	CliNonce []byte
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
	crcId := make([]byte, 0, 4)
	binary.BigEndian.PutUint32(crcId, 1)
	data := append(crcId, nonce[:]...)
	dataLen := make([]byte, 0, 4)
	binary.BigEndian.PutUint32(dataLen, uint32(len(data)))
	buf := bytes.NewBuffer(dataLen)
	buf.Write(data)
	_, err := conn.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	rsp := <-cliRead
	data = rsp[4:]
	res := &ResPqMulti{}
	if err = json.Unmarshal(data, res); err != nil {
		log.Fatal(err.Error())
		return
	}
	cliVar.cliNonce = res.Nonce
	cliVar.fingers = res.Finger
	return
}

func reqECDH() {
	pubKey, privKey := protoCli.GenECDHKey()
}
