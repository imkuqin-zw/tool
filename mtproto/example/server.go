package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
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
	cliNonce []byte
}

type ResPqMulti struct {
	Nonce  []byte
	Finger []uint64
}

func main() {
	var err error
	protoSvr, err = mtproto.NewServer("")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:30000")
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
			var conn, err = t.AcceptTCP()
			if err != nil || conn == nil {
				log.Println(err.Error())
				continue
			}
			go server(conn)
		}
	}()

	cancel()
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
				_ = conn.Close()
				log.Fatal(err.Error())
			}
			dataLen := binary.BigEndian.Uint32(lenData[:])
			data := make([]byte, 0, dataLen)
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
				n, err := conn.Write(data)
				if err != nil {
					_ = conn.Close()
					log.Fatal(err.Error())
				}
				start += n
			}
		}
	}
}

func sendData(data []byte) {
	dataLen := make([]byte, 0, 4)
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
			case 3:
			default:
				log.Println("cmd not found", cmd)
			}
		case <-ctx.Done():
			return
		}

	}
}

func reqPqMultiHandle(req []byte) {
	svrVar.cliNonce = req[:4]
	nonce := protoSvr.GenNonce128()
	res := ResPqMulti{
		Nonce:  nonce[:],
		Finger: protoSvr.GetRsaKeyFingers(),
	}
	data, _ := json.Marshal(res)
	sendData(append(req[:4], data...))
}
