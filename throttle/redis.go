package throttle

import (
	"github.com/imkuqin-zw/tool/cache"
	"github.com/garyburd/redigo/redis"
	"time"
)

type redisCatch struct{
	ip string
	prefix string
}

func NewRedisCatch(ip string) CatchThrottle {
	return redisCatch{ip: ip, prefix: "throttle_"}
}

func (r redisCatch) Attempts(key string) (val int, err error) {
	conn, err := cache.GetRedisConn(r.ip, r.prefix)
	if err != nil {
		return
	}
	defer conn.Close()
	val, err = redis.Int(conn.Get("key"))
	if err != nil && err == redis.ErrNil {
		err = conn.Set(key, 0, 0)
		return
	}
	return
}

func (r redisCatch) Has(key string) (exist bool, err error) {
	conn, err := cache.GetRedisConn(r.ip, r.prefix)
	if err != nil {
		return
	}
	defer conn.Close()
	exist, err = conn.Exist(key)
	return
}

func (r redisCatch) ResetAttempts(key string) (err error)  {
	conn, err := cache.GetRedisConn(r.ip, r.prefix)
	if err != nil {
		return
	}
	defer conn.Close()
	err = conn.Del(key)
	return
}

func (r redisCatch) Hit(key string, decayMinutes int) (hits int) {
	conn, err := cache.GetRedisConn(r.ip, r.prefix)
	if err != nil {
		return
	}
	defer conn.Close()
	availableAt := int(time.Now().Unix()) + decayMinutes * 60
	conn.SetNx(key + ":time", availableAt, decayMinutes * 60)
	if err != nil {
		return
	}
	added, err := conn.SetNx(key, 0, decayMinutes * 60)
	if err != nil {
		return
	}
	hits, err = redis.Int(conn.Increment(key))
	if err != nil {
		return
	}
	if !added && hits == 1 {
		conn.Set(key, 1, decayMinutes * 60)
	}
	return
}

func (r redisCatch) availableIn(key string) (validTime int, err error) {
	conn, err := cache.GetRedisConn(r.ip, r.prefix)
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