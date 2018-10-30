package throttle

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/imkuqin-zw/tool/cache"
	"time"
)

const LAST_TS = "last_ts"
const COUNT = "cnt"

type redisSlidingWindow struct {
	ip     string
	prefix string
}

func (r *redisSlidingWindow) DecrCount(key string, count int) (int, error) {
	conn, err := cache.NewConn(r.prefix, r.ip)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return redis.Int(conn.HIncrby(key, COUNT, -count))
}

func (r *redisSlidingWindow) GetLastTs(key string) (int64, error) {
	conn, err := cache.NewConn(r.prefix, r.ip)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	lastTs, err := redis.Int64(conn.HGet(key, LAST_TS))
	if err != nil {
		if err != redis.ErrNil {
			return 0, err
		}
	}
	return lastTs, nil
}

func (r *redisSlidingWindow) Hit(key string, curWind, now, duration int64) error {
	conn, err := cache.NewConn(r.prefix, r.ip)
	if err != nil {
		return err
	}
	defer conn.Close()
	if err = conn.HSet(key, LAST_TS, now); err != nil {
		return err
	}
	if _, err = conn.HIncrby(key, COUNT, 1); err != nil {
		return err
	}
	if _, err = conn.HIncrby(key, fmt.Sprintf("%s:%d", COUNT, curWind), 1); err != nil {
		return err
	}
	return conn.Expire(key, time.Duration(duration)*time.Second)
}

func (r *redisSlidingWindow) DelWind(key string, filed ...interface{}) (int, error) {
	conn, err := cache.NewConn(r.prefix, r.ip)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	var decr int
	decrArr, err := redis.Ints(conn.HMGet(key, filed...))
	if err != nil {
		return 0, err
	}
	if err = conn.HDel(key, filed...); err != nil {
		return 0, err
	}
	for _, item := range decrArr {
		decr += item
	}
	return decr, nil
}
