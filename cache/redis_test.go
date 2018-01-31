package cache

import (
	"testing"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func init() {
	Register("127.0.0.1:6379", nil)
}

func TestRedisConn_Increment(t *testing.T) {
	conn, err := NewConn("throttle_")
	if err != nil {
		return
	}
	defer conn.Close()

	hits, err := redis.Int(conn.Increment("dabb429575eb54a43bfed3848c701c0de87b9354"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(hits)
	return
}

func TestRedisConn_SetNx(t *testing.T) {
	conn, err := NewConn("throttle_")
	if err != nil {
		return
	}
	defer conn.Close()
	added, err := conn.SetNx("dabb429575eb54a43bfed3848c701c0de87b9354", 0, 60)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(added)
}

func TestRedisConn_Get(t *testing.T) {
	conn, err := NewConn("throttle_")
	defer conn.Close()
	if err != nil {
		return
	}
	result, err := redis.Int(conn.Get("dabb429575eb54a43bfed3848c701c0de87b9354"))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(result)
}