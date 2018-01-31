package cache

import "time"

type Conn interface{
	Close() (err error)

	Get(key string) (value interface{}, err error)

	Set(key string, val interface{}, expire int) (err error)

	Del(key string) (err error)

	Exist(key string) (exist bool, err error)

	Increment(key string) (val interface{}, err error)

	SetNx(key string, val interface{}, expire int) (value bool, err error)

	RPush(key string, val interface{}) (err error)

	LPush(key string, val interface{}) (err error)

	LPop(key string, val interface{}) (value interface{}, err error)

	RPop(key string, val interface{}) (value interface{}, err error)

	Do(command, key string, params ...interface{}) (value interface{}, err error)

	Lock(key string) (err error)

	UnLock(key string) (err error)
}

type RedisConfig struct {
	MaxIdle 		int
	MaxActive 		int
	IdleTimeOut 	time.Duration
	ConnectTimeout 	time.Duration
	ReadTimeout		time.Duration
	WriteTimeout	time.Duration
}

func Register(ip string, config *RedisConfig) {
	registerRedis(ip, config)
}

func Default(ip string) error {
	return setDefaultRedis(ip)
}

func NewConn(config ...string) (Conn, error) {
	return getRedisConn(config...)
}