package cache

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"fmt"
)

var redisPool map[string]*redis.Pool
var defaultIp string

var NoIpErr = fmt.Errorf("[Redis] ip not register")
var NoDefaultErr = fmt.Errorf("[Redis] default not register")
var DeadLockErr = fmt.Errorf("[Redis] dead lock")
var NotClose = fmt.Errorf("[Redis] not close")

type redisConn struct {
	conn redis.Conn
	prefix string
}

func init() {
	redisPool = make(map[string]*redis.Pool)
}

func registerRedis(ip string, config *RedisConfig) {
	_, ok := redisPool[ip]
	if ok {
		return
	}
	if config == nil {
		config = &RedisConfig{
			MaxIdle: 80,
			MaxActive: 12000,
			IdleTimeOut: 180 * time.Second,
			ConnectTimeout: time.Second * 10,
			ReadTimeout: time.Second * 10,
			WriteTimeout: time.Second * 10,
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
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				println(err)
			}
			return err
		},
	}
	if len(redisPool) == 1 {
		defaultIp = ip
	}
	return
}

func setDefaultRedis(ip string) error {
	_, ok := redisPool[ip]
	if !ok {
		return NoIpErr
	}
	defaultIp = ip
	return nil
}

func getRedisConn(config ...string) (Conn, error) {
	if defaultIp == "" {
		return nil, NoDefaultErr
	}
	conn := &redisConn{}
	paramLen := len(config)
	if paramLen == 0 {
		conn.conn = redisPool[defaultIp].Get()
	} else if paramLen == 1 {
		conn.conn = redisPool[defaultIp].Get()
		conn.prefix = config[0]
	} else {
		pool, ok := redisPool[config[1]]
		if !ok {
			return nil, NoIpErr
		}
		conn.conn = pool.Get()
		if paramLen > 1 {
			conn.prefix = config[0]
		}
	}
	return conn, nil
}

func (r *redisConn) Close() (err error) {
	err = r.conn.Close()
	return
}

func (r *redisConn) Get(key string) (value interface{}, err error) {
	key = r.prefix + key
	value, err = r.conn.Do("GET", key)
	return
}

func (r *redisConn) Set(key string, val interface{}, expire int) (err error) {
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

func (r *redisConn) Exist(key string) (exist bool, err error) {
	key = r.prefix + key
	exist, err = redis.Bool(r.conn.Do("EXISTS", key))
	return
}

func (r *redisConn) Increment(key string) (val interface{}, err error) {
	key = r.prefix + key
	val, err = r.conn.Do("INCR", key)
	return
}

func (r *redisConn) SetNx(key string, val interface{}, expire int) (value bool, err error) {
	key = r.prefix + key
	notExist, err := redis.Bool(r.conn.Do("SETNX", key, val))
	if err != nil || !notExist {
		return
	}
	if expire != 0 {
		_, err = r.conn.Do("EXPIRE", key, expire)
	}
	return
}

func (r *redisConn) Del(key string) (err error) {
	key = r.prefix + key
	_, err = r.conn.Do("DEL", key)
	return
}

func (r *redisConn) RPush(key string, val interface{}) (err error) {
	key = r.prefix + key
	_, err = r.conn.Do("RPUSH", key, val)
	return
}

func (r *redisConn) LPush(key string, val interface{}) (err error) {
	key = r.prefix + key
	_, err = r.conn.Do("LPUSH", key, val)
	return
}

func (r *redisConn) LPop(key string, val interface{}) (value interface{}, err error) {
	key = r.prefix + key
	value, err = r.conn.Do("LPOP", key)
	return
}

func (r *redisConn) RPop(key string, val interface{}) (value interface{}, err error) {
	key = r.prefix + key
	value, err = r.conn.Do("RPOP", key)
	return
}

func (r *redisConn) Do(command, key string, params ...interface{}) (value interface{}, err error) {
	val := []interface{}{r.prefix + key}
	val = append(val, params...)
	value, err = r.conn.Do(command, val...)
	return
}

func (r *redisConn) Lock(key string) (err error) {
	startTime := time.Now()
	for {
		var result interface{}
		result, err = r.Get(key)
		if err != nil {
			return
		}
		if result == nil {
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

func (r *redisConn) UnLock(key string) (err error) {
	err = r.Del(key)
	return
}
