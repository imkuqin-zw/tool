package fastdfs

import (
	"bytes"
	"encoding/binary"
	"os"
)

const (
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ONE = 101
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ONE               = 102
	TRACKER_PROTO_CMD_SERVICE_QUERY_UPDATE                  = 103
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ONE    = 104
	TRACKER_PROTO_CMD_SERVICE_QUERY_FETCH_ALL               = 105
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITHOUT_GROUP_ALL = 106
	TRACKER_PROTO_CMD_SERVICE_QUERY_STORE_WITH_GROUP_ALL    = 107
	TRACKER_PROTO_CMD_RESP                                  = 100
	FDFS_PROTO_CMD_ACTIVE_TEST                              = 111 //active test, tracker and storage both support since V1.28

	STORAGE_PROTO_CMD_REPORT_CLIENT_IP      = 9 //ip as tracker client
	STORAGE_PROTO_CMD_UPLOAD_FILE           = 11
	STORAGE_PROTO_CMD_DELETE_FILE           = 12
	STORAGE_PROTO_CMD_SET_METADATA          = 13
	STORAGE_PROTO_CMD_DOWNLOAD_FILE         = 14
	STORAGE_PROTO_CMD_GET_METADATA          = 15
	STORAGE_PROTO_CMD_SYNC_CREATE_FILE      = 16
	STORAGE_PROTO_CMD_SYNC_DELETE_FILE      = 17
	STORAGE_PROTO_CMD_SYNC_UPDATE_FILE      = 18
	STORAGE_PROTO_CMD_SYNC_CREATE_LINK      = 19
	STORAGE_PROTO_CMD_CREATE_LINK           = 20
	STORAGE_PROTO_CMD_UPLOAD_SLAVE_FILE     = 21
	STORAGE_PROTO_CMD_QUERY_FILE_INFO       = 22
	STORAGE_PROTO_CMD_UPLOAD_APPENDER_FILE  = 23 //create appender file
	STORAGE_PROTO_CMD_APPEND_FILE           = 24 //append file
	STORAGE_PROTO_CMD_SYNC_APPEND_FILE      = 25
	STORAGE_PROTO_CMD_FETCH_ONE_PATH_BINLOG = 26 //fetch binlog of one store path
	STORAGE_PROTO_CMD_RESP                  = TRACKER_PROTO_CMD_RESP
	STORAGE_PROTO_CMD_UPLOAD_MASTER_FILE    = STORAGE_PROTO_CMD_UPLOAD_FILE

	STORAGE_PROTO_CMD_TRUNK_ALLOC_SPACE          = 27 //since V3.00
	STORAGE_PROTO_CMD_TRUNK_ALLOC_CONFIRM        = 28 //since V3.00
	STORAGE_PROTO_CMD_TRUNK_FREE_SPACE           = 29 //since V3.00
	STORAGE_PROTO_CMD_TRUNK_SYNC_BINLOG          = 30 //since V3.00
	STORAGE_PROTO_CMD_TRUNK_GET_BINLOG_SIZE      = 31 //since V3.07
	STORAGE_PROTO_CMD_TRUNK_DELETE_BINLOG_MARKS  = 32 //since V3.07
	STORAGE_PROTO_CMD_TRUNK_TRUNCATE_BINLOG_FILE = 33 //since V3.07

	STORAGE_PROTO_CMD_MODIFY_FILE        = 34 //since V3.08
	STORAGE_PROTO_CMD_SYNC_MODIFY_FILE   = 35 //since V3.08
	STORAGE_PROTO_CMD_TRUNCATE_FILE      = 36 //since V3.08
	STORAGE_PROTO_CMD_SYNC_TRUNCATE_FILE = 37 //since V3.08

	//for overwrite all old metadata
	STORAGE_SET_METADATA_FLAG_OVERWRITE     = 'O'
	STORAGE_SET_METADATA_FLAG_OVERWRITE_STR = "O"
	//for replace, insert when the meta item not exist, otherwise update it
	STORAGE_SET_METADATA_FLAG_MERGE     = 'M'
	STORAGE_SET_METADATA_FLAG_MERGE_STR = "M"

	FDFS_RECORD_SEPERATOR = '\x01'
	FDFS_FIELD_SEPERATOR  = '\x02'

	//common constants
	FDFS_GROUP_NAME_MAX_LEN     = 16
	IP_ADDRESS_SIZE             = 16
	FDFS_PROTO_PKG_LEN_SIZE     = 8
	FDFS_PROTO_CMD_SIZE         = 1
	FDFS_PROTO_STATUS_SIZE      = 1
	FDFS_PROTO_IP_PORT_SIZE     = (IP_ADDRESS_SIZE + 6)
	FDFS_MAX_SERVERS_EACH_GROUP = 32
	FDFS_MAX_GROUPS             = 512
	FDFS_MAX_TRACKERS           = 16
	FDFS_DOMAIN_NAME_MAX_LEN    = 128

	FDFS_MAX_META_NAME_LEN  = 64
	FDFS_MAX_META_VALUE_LEN = 256

	FDFS_FILE_PREFIX_MAX_LEN    = 16
	FDFS_LOGIC_FILE_PATH_LEN    = 10
	FDFS_TRUE_FILE_PATH_LEN     = 6
	FDFS_FILENAME_BASE64_LENGTH = 27
	FDFS_TRUNK_FILE_INFO_LEN    = 16
	FDFS_FILE_EXT_NAME_MAX_LEN  = 6
	FDFS_SPACE_SIZE_BASE_INDEX  = 2 // storage space size based (MB)

	FDFS_UPLOAD_BY_BUFFER   = 1
	FDFS_UPLOAD_BY_FILENAME = 2
	FDFS_UPLOAD_BY_FILE     = 3
	FDFS_DOWNLOAD_TO_BUFFER = 1
	FDFS_DOWNLOAD_TO_FILE   = 2

	FDFS_NORMAL_LOGIC_FILENAME_LENGTH = (FDFS_LOGIC_FILE_PATH_LEN + FDFS_FILENAME_BASE64_LENGTH + FDFS_FILE_EXT_NAME_MAX_LEN + 1)

	FDFS_TRUNK_FILENAME_LENGTH       = (FDFS_TRUE_FILE_PATH_LEN + FDFS_FILENAME_BASE64_LENGTH + FDFS_TRUNK_FILE_INFO_LEN + 1 + FDFS_FILE_EXT_NAME_MAX_LEN)
	FDFS_TRUNK_LOGIC_FILENAME_LENGTH = (FDFS_TRUNK_FILENAME_LENGTH + (FDFS_LOGIC_FILE_PATH_LEN - FDFS_TRUE_FILE_PATH_LEN))

	FDFS_VERSION_SIZE = 6

	TRACKER_QUERY_STORAGE_FETCH_BODY_LEN = (FDFS_GROUP_NAME_MAX_LEN + IP_ADDRESS_SIZE - 1 + FDFS_PROTO_PKG_LEN_SIZE)
	TRACKER_QUERY_STORAGE_STORE_BODY_LEN = (FDFS_GROUP_NAME_MAX_LEN + IP_ADDRESS_SIZE - 1 + FDFS_PROTO_PKG_LEN_SIZE + 1)
	//status code, order is important!
	FDFS_STORAGE_STATUS_INIT       = 0
	FDFS_STORAGE_STATUS_WAIT_SYNC  = 1
	FDFS_STORAGE_STATUS_SYNCING    = 2
	FDFS_STORAGE_STATUS_IP_CHANGED = 3
	FDFS_STORAGE_STATUS_DELETED    = 4
	FDFS_STORAGE_STATUS_OFFLINE    = 5
	FDFS_STORAGE_STATUS_ONLINE     = 6
	FDFS_STORAGE_STATUS_ACTIVE     = 7
	FDFS_STORAGE_STATUS_RECOVERY   = 9
	FDFS_STORAGE_STATUS_NONE       = 99
)

type StorageSvr struct {
	ipAddr         string
	port           int
	groupName      string
	storePathIndex int
}

type uploadFileRequest struct {
	storePathIndex uint8
	fileSize       int64
	fileExtName    string
}

func (this *uploadFileRequest) marshal() []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteByte(byte(this.storePathIndex))
	binary.Write(buffer, binary.BigEndian, this.fileSize)

	// 6 bit fileExtName
	fileExtNameBytes := bytes.NewBufferString(this.fileExtName).Bytes()
	for i := 0; i < 6; i++ {
		if i >= len(fileExtNameBytes) {
			buffer.WriteByte(byte(0))
		} else {
			buffer.WriteByte(fileExtNameBytes[i])
		}
	}
	return buffer.Bytes()
}

type uploadSlaveFileRequest struct {
	masterFilenameLen int64
	fileSize          int64
	prefixName        string
	fileExtName       string
	masterFilename    string
}

// #slave_fmt |-master_len(8)-file_size(8)-prefix_name(16)-file_ext_name(6)
// #           -master_name(master_filename_len)-|
func (this *uploadSlaveFileRequest) marshal() []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, this.masterFilenameLen)
	binary.Write(buffer, binary.BigEndian, this.fileSize)

	// 16 bit prefixName
	prefixNameBytes := bytes.NewBufferString(this.prefixName).Bytes()
	for i := 0; i < 16; i++ {
		if i >= len(prefixNameBytes) {
			buffer.WriteByte(byte(0))
		} else {
			buffer.WriteByte(prefixNameBytes[i])
		}
	}

	// 6 bit fileExtName
	fileExtNameBytes := bytes.NewBufferString(this.fileExtName).Bytes()
	for i := 0; i < 6; i++ {
		if i >= len(fileExtNameBytes) {
			buffer.WriteByte(byte(0))
		} else {
			buffer.WriteByte(fileExtNameBytes[i])
		}
	}

	// master_filename_len bit master_name
	masterFilenameBytes := bytes.NewBufferString(this.masterFilename).Bytes()
	for i := 0; i < int(this.masterFilenameLen); i++ {
		buffer.WriteByte(masterFilenameBytes[i])
	}
	return buffer.Bytes()
}

type UploadFileResp struct {
	GroupName    string
	RemoteFileId string
}

// recv_fmt: |-group_name(16)-remote_file_name(recv_size - 16)-|
func (this *UploadFileResp) unmarshal(data []byte) error {
	buff := bytes.NewBuffer(data)
	var err error
	this.GroupName, err = readCstr(buff, FDFS_GROUP_NAME_MAX_LEN)
	if err != nil {
		return err
	}
	remoteFilename := string(data[len(data)-buff.Len():])
	this.RemoteFileId = this.GroupName + string(os.PathSeparator) + remoteFilename
	return nil
}

type deleteFileRequest struct {
	groupName      string
	remoteFilename string
}

// #del_fmt: |-group_name(16)-filename(len)-|
func (this *deleteFileRequest) marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)

	// 16 bit groupName
	groupNameBytes := bytes.NewBufferString(this.groupName).Bytes()
	for i := 0; i < 16; i++ {
		if i >= len(groupNameBytes) {
			buffer.WriteByte(byte(0))
		} else {
			buffer.WriteByte(groupNameBytes[i])
		}
	}

	// remoteFilenameLen bit remoteFilename
	remoteFilenameBytes := bytes.NewBufferString(this.remoteFilename).Bytes()
	for i := 0; i < len(remoteFilenameBytes); i++ {
		buffer.WriteByte(remoteFilenameBytes[i])
	}
	return buffer.Bytes(), nil
}

type DeleteFileResponse struct {
	groupName      string
	remoteFilename string
}

// recv_fmt: |-group_name(16)-remote_file_name(recv_size - 16)-|
func (this *DeleteFileResponse) unmarshal(data []byte) error {
	buff := bytes.NewBuffer(data)
	var err error
	this.groupName, err = readCstr(buff, FDFS_GROUP_NAME_MAX_LEN)
	if err != nil {
		return err
	}
	this.remoteFilename = string(data[len(data)-buff.Len():])
	return nil
}

type downloadFileRequest struct {
	offset         int64
	downloadSize   int64
	groupName      string
	remoteFilename string
}

// #down_fmt: |-offset(8)-download_bytes(8)-group_name(16)-remote_filename(len)-|
func (this *downloadFileRequest) marshal() ([]byte, error) {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, this.offset)
	binary.Write(buffer, binary.BigEndian, this.downloadSize)

	// 16 bit groupName
	groupNameBytes := bytes.NewBufferString(this.groupName).Bytes()
	for i := 0; i < 16; i++ {
		if i >= len(groupNameBytes) {
			buffer.WriteByte(byte(0))
		} else {
			buffer.WriteByte(groupNameBytes[i])
		}
	}

	// remoteFilenameLen bit remoteFilename
	remoteFilenameBytes := bytes.NewBufferString(this.remoteFilename).Bytes()
	for i := 0; i < len(remoteFilenameBytes); i++ {
		buffer.WriteByte(remoteFilenameBytes[i])
	}
	return buffer.Bytes(), nil
}

type DownloadFileResponse struct {
	RemoteFileId string
	Content      interface{}
	DownloadSize int64
}
