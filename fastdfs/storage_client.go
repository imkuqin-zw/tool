package fastdfs

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os"
)

type StorageClient struct {
	pool Pool
}

func (this *StorageClient) uploadByFilename(storeSvr *StorageSvr, filename string) (*UploadFileResp, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	fileExtName := getFileExt(filename)

	return this.uploadFile(storeSvr, filename, int64(fileSize), FDFS_UPLOAD_BY_FILENAME,
		STORAGE_PROTO_CMD_UPLOAD_FILE, "", "", fileExtName)
}

func (this *StorageClient) uploadByBuffer(storeSvr *StorageSvr, fileBuffer []byte,
	fileExtName string) (*UploadFileResp, error) {
	return this.uploadFile(storeSvr, fileBuffer, int64(len(fileBuffer)), FDFS_UPLOAD_BY_BUFFER,
		STORAGE_PROTO_CMD_UPLOAD_FILE, "", "", fileExtName)
}

func (this *StorageClient) uploadSlaveByFilename(storeSvr *StorageSvr, filename string,
	prefixName string, remoteFileId string) (*UploadFileResp, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	fileExtName := getFileExt(filename)

	return this.uploadFile(storeSvr, filename, int64(fileSize), FDFS_UPLOAD_BY_FILENAME,
		STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE, remoteFileId, prefixName, fileExtName)
}

func (this *StorageClient) uploadSlaveByBuffer(storeSvr *StorageSvr, fileBuffer []byte,
	remoteFileId string, fileExtName string) (*UploadFileResp, error) {

	return this.uploadFile(storeSvr, fileBuffer, int64(len(fileBuffer)), FDFS_UPLOAD_BY_BUFFER,
		STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE, "", remoteFileId, fileExtName)
}

func (this *StorageClient) uploadAppenderByFilename(storeSvr *StorageSvr, filename string) (*UploadFileResp, error) {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	fileExtName := getFileExt(filename)

	return this.uploadFile(storeSvr, filename, int64(fileSize), FDFS_UPLOAD_BY_FILENAME,
		STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE, "", "", fileExtName)
}

func (this *StorageClient) uploadAppenderByBuffer(storeSvr *StorageSvr, fileBuffer []byte,
	fileExtName string) (*UploadFileResp, error) {
	bufferSize := len(fileBuffer)

	return this.uploadFile(storeSvr, fileBuffer, int64(bufferSize), FDFS_UPLOAD_BY_BUFFER,
		STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE, "", "", fileExtName)
}

func (this *StorageClient) uploadFile(storeSvr *StorageSvr, fileContent interface{}, fileSize int64,
	uploadType int, cmd int8, masterFilename string, prefixName string, fileExtName string) (*UploadFileResp, error) {

	var (
		uploadSlave bool
		headerLen   int64 = 15
		reqBuf      []byte
	)
	connNode, err := this.pool.Get()
	if err != nil {
		return nil, err
	}
	conn := connNode.c

	//request header
	masterFilenameLen := int64(len(masterFilename))
	if len(storeSvr.groupName) > 0 && len(masterFilename) > 0 {
		uploadSlave = true
		// #slave_fmt |-master_len(8)-file_size(8)-prefix_name(16)-file_ext_name(6)
		//       #           -master_name(master_filename_len)-|
		headerLen = int64(38) + masterFilenameLen
	}
	h := &header{
		pkgLen: headerLen + int64(fileSize),
		cmd:    cmd,
	}
	if err := conn.WriteHeader(h); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("write header fault", zap.Error(err))
		return nil, err
	}

	// request body
	if uploadSlave {
		req := &uploadSlaveFileRequest{}
		req.masterFilenameLen = masterFilenameLen
		req.fileSize = int64(fileSize)
		req.prefixName = prefixName
		req.fileExtName = fileExtName
		req.masterFilename = masterFilename
		reqBuf = req.marshal()
	} else {
		req := &uploadFileRequest{}
		req.storePathIndex = uint8(storeSvr.storePathIndex)
		req.fileSize = int64(fileSize)
		req.fileExtName = fileExtName
		reqBuf = req.marshal()
	}
	if err = conn.Write(reqBuf); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("write body fault", zap.Error(err))
		return nil, err
	}

	var filebuf []byte
	switch uploadType {
	case FDFS_UPLOAD_BY_FILENAME:
		if filename, ok := fileContent.(string); ok {
			filebuf, err = GetFileData(filename)
			if err != nil {
				this.pool.Put(connNode, true)
				Logger.Error("get file data fault", zap.Error(err))
				return nil, err
			}
		} else {
			this.pool.Put(connNode, true)
			Logger.Error("fileContent type fault")
			return nil, errors.New("fileContent type fault")
		}
	case FDFS_UPLOAD_BY_BUFFER:
		if fileBuffer, ok := fileContent.([]byte); ok {
			filebuf = fileBuffer
		} else {
			this.pool.Put(connNode, true)
			Logger.Error("fileContent type fault")
			return nil, errors.New("fileContent type fault")
		}
	default:
		this.pool.Put(connNode, true)
		Logger.Error("upload type not be defined")
		return nil, errors.New("upload type not be defined")
	}
	if err = conn.Write(filebuf); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("write file data fault", zap.Error(err))
		return nil, err
	}

	// respose header
	if h, err = conn.ReadHeader(); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("read header fault", zap.Error(err))
		return nil, err
	}
	if h.status != 0 {
		this.pool.Put(connNode, false)
		return nil, Errno{int(h.status)}
	}

	//　response body
	rspBuf, err := conn.ReadN(h.pkgLen)
	if err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("read body fault", zap.Error(err))
		return nil, err
	}
	this.pool.Put(connNode, false)
	ur := &UploadFileResp{}
	if err = ur.unmarshal(rspBuf); err != nil {
		Logger.Error("recvBuf can not unmarshal", zap.Error(err))
		return nil, err
	}
	return ur, nil
}

func (this *StorageClient) appendByFilename() {

}

func (this *StorageClient) appendByBuffer() {

}

func (this *StorageSvr) appendFile(groupName string, appenderFileName string, fileSize int64) {

}

func (this *StorageClient) storageDeleteFile(tc *TrackerClient, storeServ *StorageServer, remoteFilename string) error {
	var (
		conn   net.Conn
		reqBuf []byte
		err    error
	)

	conn, err = this.pool.Get()
	defer conn.Close()
	if err != nil {
		return err
	}

	th := &trackerHeader{}
	th.cmd = STORAGE_PROTO_CMD_DELETE_FILE
	fileNameLen := len(remoteFilename)
	th.pkgLen = int64(FDFS_GROUP_NAME_MAX_LEN + fileNameLen)
	th.sendHeader(conn)

	req := &deleteFileRequest{}
	req.groupName = storeServ.groupName
	req.remoteFilename = remoteFilename
	reqBuf, err = req.marshal()
	if err != nil {
		logger.Warnf("deleteFileRequest.marshal error :%s", err.Error())
		return err
	}
	TcpSendData(conn, reqBuf)

	th.recvHeader(conn)
	if th.status != 0 {
		return Errno{int(th.status)}
	}
	return nil
}

func (this *StorageClient) storageDownloadToFile(tc *TrackerClient,
	storeServ *StorageServer, localFilename string, offset int64,
	downloadSize int64, remoteFilename string) (*DownloadFileResponse, error) {
	return this.storageDownloadFile(tc, storeServ, localFilename, offset, downloadSize, FDFS_DOWNLOAD_TO_FILE, remoteFilename)
}

func (this *StorageClient) storageDownloadToBuffer(tc *TrackerClient,
	storeServ *StorageServer, fileBuffer []byte, offset int64,
	downloadSize int64, remoteFilename string) (*DownloadFileResponse, error) {
	return this.storageDownloadFile(tc, storeServ, fileBuffer, offset, downloadSize, FDFS_DOWNLOAD_TO_BUFFER, remoteFilename)
}

func (this *StorageClient) storageDownloadFile(tc *TrackerClient,
	storeServ *StorageServer, fileContent interface{}, offset int64, downloadSize int64,
	downloadType int, remoteFilename string) (*DownloadFileResponse, error) {

	var (
		conn          net.Conn
		reqBuf        []byte
		localFilename string
		recvBuff      []byte
		recvSize      int64
		err           error
	)

	conn, err = this.pool.Get()
	defer conn.Close()
	if err != nil {
		return nil, err
	}

	th := &trackerHeader{}
	th.cmd = STORAGE_PROTO_CMD_DOWNLOAD_FILE
	th.pkgLen = int64(FDFS_PROTO_PKG_LEN_SIZE*2 + FDFS_GROUP_NAME_MAX_LEN + len(remoteFilename))
	th.sendHeader(conn)

	req := &downloadFileRequest{}
	req.offset = offset
	req.downloadSize = downloadSize
	req.groupName = storeServ.groupName
	req.remoteFilename = remoteFilename
	reqBuf, err = req.marshal()
	if err != nil {
		logger.Warnf("downloadFileRequest.marshal error :%s", err.Error())
		return nil, err
	}
	TcpSendData(conn, reqBuf)

	th.recvHeader(conn)
	if th.status != 0 {
		return nil, Errno{int(th.status)}
	}

	switch downloadType {
	case FDFS_DOWNLOAD_TO_FILE:
		if localFilename, ok := fileContent.(string); ok {
			recvSize, err = TcpRecvFile(conn, localFilename, th.pkgLen)
		}
	case FDFS_DOWNLOAD_TO_BUFFER:
		if _, ok := fileContent.([]byte); ok {
			recvBuff, recvSize, err = TcpRecvResponse(conn, th.pkgLen)
		}
	}
	if err != nil {
		logger.Warnf(err.Error())
		return nil, err
	}
	if recvSize < downloadSize {
		errmsg := "[-] Error: Storage response length is not match, "
		errmsg += fmt.Sprintf("expect: %d, actual: %d", th.pkgLen, recvSize)
		logger.Warn(errmsg)
		return nil, errors.New(errmsg)
	}

	dr := &DownloadFileResponse{}
	dr.RemoteFileId = storeServ.groupName + string(os.PathSeparator) + remoteFilename
	if downloadType == FDFS_DOWNLOAD_TO_FILE {
		dr.Content = localFilename
	} else {
		dr.Content = recvBuff
	}
	dr.DownloadSize = recvSize
	return dr, nil
}
