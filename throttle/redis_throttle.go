package throttle

import (
	"github.com/garyburd/redigo/redis"
	"github.com/imkuqin-zw/tool/cache"
	"time"
)

type redisThrottle struct {
	ip     string
	prefix string
}

func NewRedisCatch(ip string) CatchThrottle {
	return &redisThrottle{ip: ip, prefix: "throttle_"}
}

func (r *redisThrottle) Attempts(key string) (val int, err error) {
	conn, err := cache.NewConn(r.prefix, r.ip)
	if err != nil {
		return
	}
	defer conn.Close()
	val, err = redis.Int(conn.Get(key))
	if err != nil && err == redis.ErrNil {
		err = conn.Set(key, 0, 0)
		return
	}
	return
}

func (r *redisThrottle) Has(key string) (exist bool, err error) {
	conn, err := cache.NewConn(r.prefix, r.ip)
	if err != nil {
		return
	}
	defer conn.Close()
	exist, err = conn.Exist(key)
	return
}

func (r *redisThrottle) ResetAttempts(key string) (err error) {
	conn, err := cache.NewConn(r.prefix, r.ip)
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.Del(key)
	return
}

func (r *redisThrottle) Hit(key string, decayMinutes int) (hits int, err error) {
	conn, err := cache.NewConn(r.prefix, r.ip)
	if err != nil {
		return
	}
	defer conn.Close()
	availableAt := int(time.Now().Unix()) + decayMinutes*60
	_, err = conn.SetNx(key+":timer", availableAt, decayMinutes*60)
	if err != nil {
		return
	}
	added, err := conn.SetNx(key, 0, decayMinutes*60)
	if err != nil {
		return
	}

	hits, err = redis.Int(conn.Increment(key))
	if err != nil {
		return
	}
	if !added && hits == 1 {
		conn.Set(key, 1, decayMinutes*60)
	}
	return
}

func (r *redisThrottle) availableIn(key string) (validTime int, err error) {
	conn, err := cache.NewConn(r.prefix, r.ip)
	if err != nil {
		return
	}
	defer conn.Close()
	validTime, err = redis.Int(conn.Get(key + ":timer"))
	if err != nil {
		return
	}
	validTime -= int(time.Now().Unix())
	return
}
