package cache

import "time"

type Conn interface {
	Close() (err error)

	Get(key string) (value interface{}, err error)

	Set(key string, val interface{}, expire int) error

	Del(key string) (err error)

	Exist(key string) (bool, error)

	Increment(key string) (interface{}, error)

	SetNx(key string, val interface{}, expire int) (bool, error)

	RPush(key string, val interface{}) (err error)

	LPush(key string, val interface{}) (err error)

	LPop(key string, val interface{}) (interface{}, error)

	RPop(key string, val interface{}) (interface{}, error)

	Do(command, key string, params ...interface{}) (interface{}, error)

	Lock(key string) (err error)

	UnLock(key string) (err error)

	HGet(key, field string) (interface{}, error)

	HSet(key, field string, val interface{}) error

	HMGet(key string, field ...interface{}) (interface{}, error)

	HDel(key string, field ...interface{}) error

	HIncrby(key, field string, val interface{}) (interface{}, error)

	Expire(key string, val time.Duration) error
}

type RedisConfig struct {
	MaxIdle        int
	MaxActive      int
	IdleTimeOut    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
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
