package fastdfs

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type Errno struct {
	status int
}

func (e Errno) Error() string {
	errmsg := fmt.Sprintf("errno [%d] ", e.status)
	switch e.status {
	case 17:
		errmsg += "File Exist"
	case 22:
		errmsg += "Argument Invlid"
	}
	return errmsg
}

func readCstr(buff io.Reader, length int) (string, error) {
	str := make([]byte, length)
	n, err := buff.Read(str)
	if err != nil || n != len(str) {
		return "", Errno{255}
	}

	for i, v := range str {
		if v == 0 {
			str = str[0:i]
			break
		}
	}
	return string(str), nil
}
func getFileExt(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}

func splitRemoteFileId(remoteFileId string) ([]string, error) {
	parts := strings.SplitN(remoteFileId, "/", 2)
	if len(parts) < 2 {
		return nil, errors.New("error remoteFileId")
	}
	return parts, nil
}

func TcpSendData(conn net.Conn, bytesStream []byte) error {
	if _, err := conn.Write(bytesStream); err != nil {
		return err
	}
	return nil
}

func TcpSendFile(conn net.Conn, filename string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	var fileSize int64 = 0
	if fileInfo, err := file.Stat(); err == nil {
		fileSize = fileInfo.Size()
	}

	if fileSize == 0 {
		errmsg := fmt.Sprintf("file size is zeor [%s]", filename)
		return errors.New(errmsg)
	}

	fileBuffer := make([]byte, fileSize)

	_, err = file.Read(fileBuffer)
	if err != nil {
		return err
	}

	return TcpSendData(conn, fileBuffer)
}

func TcpRecvResponse(conn net.Conn, bufferSize int64) ([]byte, int64, error) {
	recvBuff := make([]byte, 0, bufferSize)
	tmp := make([]byte, 256)
	var total int64
	for {
		n, err := conn.Read(tmp)
		total += int64(n)
		recvBuff = append(recvBuff, tmp[:n]...)
		if err != nil {
			if err != io.EOF {
				return nil, 0, err
			}
			break
		}
		if total == bufferSize {
			break
		}
	}
	return recvBuff, total, nil
}

func TcpRecvFile(conn net.Conn, localFilename string, bufferSize int64) (int64, error) {
	file, err := os.Create(localFilename)
	defer file.Close()
	if err != nil {
		return 0, err
	}

	recvBuff, total, err := TcpRecvResponse(conn, bufferSize)
	if _, err := file.Write(recvBuff); err != nil {
		return 0, err
	}
	return total, nil
}
