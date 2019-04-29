package fastdfs

import (
	"bytes"
	"encoding/binary"
	"go.uber.org/zap"
)

type TrackerClient struct {
	pool Pool
}

func (this *TrackerClient) trackerQueryStorageStoreWithoutGroup() (*StorageSvr, error) {
	connNode, err := this.pool.Get()
	if err != nil {
		Logger.Error("get conn fault", zap.Error(err))
		return nil, err
	}
	conn := connNode.c

	// request header
	h := &header{
		cmd: TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE,
	}
	if err := conn.WriteHeader(h); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("write header fault", zap.Error(err))
		return nil, err
	}

	// response header
	if h, err = conn.ReadHeader(); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("read header fault", zap.Error(err))
		return nil, err
	}
	if h.status != 0 {
		this.pool.Put(connNode, false)
		return nil, Errno{int(h.status)}
	}

	// response body
	rspBuf, err := conn.ReadN(h.pkgLen)
	if err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("read body fault", zap.Error(err))
		return nil, err
	}
	this.pool.Put(connNode, false)
	svr := &StorageSvr{}
	buff := bytes.NewBuffer(rspBuf)
	svr.groupName, _ = readCstr(buff, FDFS_GROUP_NAME_MAX_LEN)
	svr.ipAddr, _ = readCstr(buff, IP_ADDRESS_SIZE-1)
	binary.Read(buff, binary.BigEndian, &svr.port)
	binary.Read(buff, binary.BigEndian, &svr.storePathIndex)
	return svr, nil
}

func (this *TrackerClient) trackerQueryStorageStorWithGroup(groupName string) (*StorageSvr, error) {
	connNode, err := this.pool.Get()
	if err != nil {
		Logger.Error("get conn fault", zap.Error(err))
		return nil, err
	}
	conn := connNode.c

	// request header
	h := &header{
		cmd:    TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE,
		pkgLen: int64(FDFS_GROUP_NAME_MAX_LEN),
	}
	if err := conn.WriteHeader(h); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("write header fault", zap.Error(err))
		return nil, err
	}

	// request body
	groupBuffer := new(bytes.Buffer)
	// 16 bit groupName
	groupNameBytes := bytes.NewBufferString(groupName).Bytes()
	for i := 0; i < 16; i++ {
		if i >= len(groupNameBytes) {
			groupBuffer.WriteByte(byte(0))
		} else {
			groupBuffer.WriteByte(groupNameBytes[i])
		}
	}
	groupBytes := groupBuffer.Bytes()
	if err = conn.Write(groupBytes); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("write body fault", zap.Error(err))
		return nil, err
	}

	// response header
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
	svr := &StorageSvr{}
	buff := bytes.NewBuffer(rspBuf)
	svr.groupName, _ = readCstr(buff, FDFS_GROUP_NAME_MAX_LEN)
	svr.ipAddr, _ = readCstr(buff, IP_ADDRESS_SIZE-1)
	binary.Read(buff, binary.BigEndian, &svr.port)
	binary.Read(buff, binary.BigEndian, &svr.storePathIndex)
	return svr, nil
}

func (this *TrackerClient) trackerQueryStorageUpdate(groupName string, remoteFilename string) (*StorageSvr, error) {
	return this.trackerQueryStorage(groupName, remoteFilename, TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE)
}

func (this *TrackerClient) trackerQueryStorageFetch(groupName string, remoteFilename string) (*StorageSvr, error) {
	return this.trackerQueryStorage(groupName, remoteFilename, TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE)
}

func (this *TrackerClient) trackerQueryStorage(groupName string, remoteFilename string, cmd int8) (*StorageSvr, error) {
	connNode, err := this.pool.Get()
	if err != nil {
		Logger.Error("get conn fault", zap.Error(err))
		return nil, err
	}
	conn := connNode.c

	// request header
	h := &header{
		cmd:    cmd,
		pkgLen: int64(FDFS_GROUP_NAME_MAX_LEN + len(remoteFilename)),
	}
	if err := conn.WriteHeader(h); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("write header fault", zap.Error(err))
		return nil, err
	}

	// #query_fmt: |-group_name(16)-filename(file_name_len)-|
	queryBuffer := new(bytes.Buffer)
	// 16 bit groupName
	groupNameBytes := bytes.NewBufferString(groupName).Bytes()
	for i := 0; i < 16; i++ {
		if i >= len(groupNameBytes) {
			queryBuffer.WriteByte(byte(0))
		} else {
			queryBuffer.WriteByte(groupNameBytes[i])
		}
	}
	// remoteFilenameLen bit remoteFilename
	remoteFilenameBytes := bytes.NewBufferString(remoteFilename).Bytes()
	for i := 0; i < len(remoteFilenameBytes); i++ {
		queryBuffer.WriteByte(remoteFilenameBytes[i])
	}
	if err = conn.Write(queryBuffer.Bytes()); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("write body fault", zap.Error(err))
		return nil, err
	}

	// response header
	if h, err = conn.ReadHeader(); err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("read header fault", zap.Error(err))
		return nil, err
	}
	if h.status != 0 {
		this.pool.Put(connNode, false)
		return nil, Errno{int(h.status)}
	}

	// response body
	rspBuf, err := conn.ReadN(h.pkgLen)
	if err != nil {
		this.pool.Put(connNode, true)
		Logger.Error("read body fault", zap.Error(err))
		return nil, err
	}
	this.pool.Put(connNode, false)
	svr := &StorageSvr{}
	buff := bytes.NewBuffer(rspBuf)
	// #recv_fmt |-group_name(16)-ipaddr(16-1)-port(8)-store_path_index(1)|
	svr.groupName, err = readCstr(buff, FDFS_GROUP_NAME_MAX_LEN)
	svr.ipAddr, err = readCstr(buff, IP_ADDRESS_SIZE-1)
	binary.Read(buff, binary.BigEndian, &svr.port)
	binary.Read(buff, binary.BigEndian, &svr.storePathIndex)
	return svr, nil
}
