package cache

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
)

type RedisConn struct {
	conn redis.Conn
	prefix string
}

var redisPool map[string]*redis.Pool
var defaultIp string

var NoIpErr = fmt.Errorf("[Redis] ip not register")
var NoDefaultErr = fmt.Errorf("[Redis] default not register")
var DeadLockErr = fmt.Errorf("[Redis] dead lock")
var NotClose = fmt.Errorf("[Redis] not close")

type RedisConfig struct {
	MaxIdle 	int
	MaxActive 	int
	IdleTimeOut time.Duration
}

func RegisterRedis(ip string, config *RedisConfig) {
	_, ok := redisPool[ip]
	if ok {
		return
	}
	if config == nil {
		config = &RedisConfig{
			MaxIdle: 80,
			MaxActive: 12000,
			IdleTimeOut: 180 * time.Second,
		}
	}
	redisPool[ip] = &redis.Pool{
		MaxIdle: config.MaxIdle,
		MaxActive: config.MaxActive,
		IdleTimeout: config.IdleTimeOut,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ip)
			return c, err
		},
	}
	if len(redisPool) == 1 {
		defaultIp = ip
	}
	return
}

func SetDefaultRedis(ip string) error {
	_, ok := redisPool[ip]
	if !ok {
		return NoIpErr
	}
	defaultIp = ip
	return nil
}

func GetRedisConn(config ...string) (redisConn *RedisConn, err error) {
	if defaultIp == "" {
		err = NoDefaultErr
		return
	}
	redisConn = &RedisConn{}
	paramLen := len(config)
	if paramLen == 0 {
		redisConn.conn = redisPool[defaultIp].Get()
	} else {
		pool, ok := redisPool[config[0]]
		if !ok {
			err = NoIpErr
			return
		}
		redisConn.conn = pool.Get()
		if paramLen > 1 {
			redisConn.prefix = config[1]
		}
	}
	return
}

func (r *RedisConn) Close() (err error) {
	err = r.conn.Close()
	return
}

func (r *RedisConn) Get(key string) (value interface{}, err error) {
	key = r.prefix + key
	value, err = r.conn.Do("GET", key)
	return
}

func (r *RedisConn) Set(key string, val interface{}, expire int) (err error) {
	key = r.prefix + key
	_, err = r.conn.Do("SET", key, val)
	if err != nil {
		return
	}
	if expire != 0 {
		_, err = r.conn.Do("EXPIRE", key, expire)
	}
	return
}

func (r *RedisConn) Exist(key string) (exist bool, err error) {
	key = r.prefix + key
	exist, err = redis.Bool(r.conn.Do("EXISTS", key))
	return
}

func (r *RedisConn) Increment(key string) (val interface{}, err error) {
	key = r.prefix + key
	val, err = r.conn.Do("INCR", key)
	return
}

func (r *RedisConn) SetNx(key string, val interface{}, expire int) (value bool, err error) {
	key = r.prefix + key
	val, err = redis.Bool(r.conn.Do("SETNX", key, val))
	if err != nil {
		return
	}
	if expire != 0 {
		_, err = r.conn.Do("EXPIRE", key, expire)
	}
	return
}

func (r *RedisConn) Del(key string) (err error) {
	key = r.prefix + key
	_, err = r.conn.Do("DEL", key)
	return
}

func (r *RedisConn) RPush(key string, val interface{}) (err error) {
	key = r.prefix + key
	_, err = r.conn.Do("RPUSH", key, val)
	return
}

func (r *RedisConn) LPush(key string, val interface{}) (err error) {
	key = r.prefix + key
	_, err = r.conn.Do("LPUSH", key, val)
	return
}

func (r *RedisConn) LPop(key string, val interface{}) (value interface{}, err error) {
	key = r.prefix + key
	value, err = r.conn.Do("LPOP", key)
	return
}

func (r *RedisConn) RPop(key string, val interface{}) (value interface{}, err error) {
	key = r.prefix + key
	value, err = r.conn.Do("RPOP", key)
	return
}

func (r *RedisConn) Do(command, key string, params ...interface{}) (value interface{}, err error) {
	val := []interface{}{r.prefix + key}
	val = append(val, params...)
	value, err = r.conn.Do(command, val...)
	return
}

func (r *RedisConn) Lock(key string) (err error) {
	startTime := time.Now()
	for {
		_, err := r.Get(key)
		if err != nil && err != redis.ErrNil {
			return
		}
		if err == redis.ErrNil {
			err = r.Set(key, "lock", 10)
			return
		}
		if time.Now().Sub(startTime) > time.Second * 4 {
			r.Del(key)
			return DeadLockErr
		}
		time.Sleep(time.Millisecond * 10)
	}
	return
}

func (r *RedisConn) UnLock(key string) (err error) {
	r.Del(key)
	return
}
