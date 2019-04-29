package fastdfs

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

const (
	HEADER_LEN = 10
)

type header struct {
	cmd    int8
	status int8
	pkgLen int64
}

type conn struct {
	c net.Conn

	readTimeout time.Duration
	r           bufio.Reader

	writeTimeout time.Duration
	w            bufio.Writer
}

func (c *conn) Close() error {
	return c.c.Close()
}

func (c *conn) WriteHeader(h *header) error {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, h.pkgLen)
	buffer.WriteByte(byte(h.cmd))
	buffer.WriteByte(byte(h.status))
	return c.Write(buffer.Bytes())
}

func (c *conn) ReadHeader() (*header, error) {
	pkgLen, err := c.ReadN(8)
	if err != nil {
		return nil, err
	}
	cmd, err := c.r.ReadByte()
	if err != nil {
		return nil, err
	}
	status, err := c.r.ReadByte()
	if err != nil {
		return nil, err
	}
	header := &header{
		pkgLen: int64(binary.BigEndian.Uint64(pkgLen)),
		cmd:    int8(cmd),
		status: int8(status),
	}
	return header, nil
}

func (c *conn) Write(data []byte) error {
	total := len(data)
	writeNum := 0
	for writeNum < total {
		nn, err := c.w.Write(data[writeNum:])
		if err != nil {
			return err
		}
		writeNum += nn
	}
	return nil
}

func (c *conn) ReadN(total int64) ([]byte, error) {
	rsp := make([]byte, total)
	var readNum int64
	for readNum < total {
		buf := make([]byte, total-readNum)
		n, err := c.r.Read(rsp)
		if err != nil {
			return nil, err
		}
		copy(rsp, buf[0:n])
		readNum += int64(n)
	}
	return rsp, nil
}
