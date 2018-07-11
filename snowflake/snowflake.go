package snowflake

import (
	"time"
	"fmt"
	"crypto/md5"
	"strconv"
	"sync/atomic"
)

var workId int64
var maxWorkId = int64(-1 ^ (-1 << 5))
var maxCenterId = int64(-1 ^ (-1 << 5))
var dataCenterId int64
var sequence int64
var lastMillisecond int64

// GetUID 得到全局唯一ID int64类型
// 首位0(1位) + 毫秒时间戳(41位) + 数据中心标识(5位) + 工作机器标识(5位) + 自增id(12位)
// 时间可以保证400年不重复
// 数据中心和机器标识一起标识节点，最多支持1024个节点
// 每个节点每一毫秒能生成最多4096个id
// 63    62          21          16      11        0
// +-----+-----------+-----------+------－+--------+
// |未使用|毫秒级时间戳 |数据中心标识 |工作机器 | 自增id  |
// +-----+-----------+-----------+------－+--------+
func GetUUID() int64 {
	now := time.Now().UnixNano() / int64(time.Millisecond)
	if now != lastMillisecond {
		sequence = 0
	} else {
		atomic.AddInt64(&sequence, 1)
	}
	lastMillisecond = now
	return now << 22 | dataCenterId << 17 | workId << 12 | sequence
}

//初始化节点标识
func Init(dataCenter int64, work int64) error {
	if maxWorkId < work {
		return fmt.Errorf("[snowflake] work_id must less than %v", maxWorkId)
	}
	if maxCenterId < dataCenter {
		return fmt.Errorf("[snowflake] work_id must less than %v", maxWorkId)
	}
	workId = work
	dataCenterId = dataCenter
	return nil
}

//获取token
func GetNewToken() string {
	uuid := strconv.FormatInt(GetUUID(), 10)
	token := md5.Sum([]byte(uuid))
	return fmt.Sprintf("%x", token)
}