package codis_cli

import (
	"fmt"
	etcdV2 "github.com/coreos/etcd/client"
	"github.com/garyburd/redigo/redis"
	"log"
	"strings"
	"testing"
	"time"
)

var client *Client
var etcdClient etcdV2.Client

func init() {
	var err error
	config := Config{
		MaxIdle:        30,
		MaxActive:      200,
		IdleTimeOut:    180 * time.Second,
		ConnectTimeout: time.Second * 10,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		RegisterKey:    "/codis3/codis-demo/proxy/",
	}
	etcdClient, err = etcdV2.New(etcdV2.Config{
		Endpoints:               strings.Split("http://192.168.2.118:2379", ","),
		Transport:               etcdV2.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second * 5,
	})
	if err != nil {
		log.Fatal("etcd init fault", err.Error())
	}
	client = NewClient(config, etcdClient)
	if err = client.Init(); err != nil {
		log.Fatal("client init fault", err.Error())
	}
}

func TestClient_Set(t *testing.T) {
	conn := client.GetConn()
	if conn == nil {
		fmt.Println("not client")
		return
	}
	defer conn.Close()
	_, err := conn.Do("set", "testfdsaf_rertr", "123")
	if err != nil {
		log.Print("error", err.Error())
		return
	}
}
func TestClient_Get(t *testing.T) {
	for {
		conn := client.GetConn()
		if conn == nil {
			fmt.Println("not client")
			return
		}
		res, err := redis.String(conn.Do("get", "testfdsaf_rertr"))
		conn.Close()
		if err != nil {
			log.Print("error", err.Error())
			return
		}
		fmt.Println(res)
		time.Sleep(100 * time.Millisecond)
	}

	close := make(chan bool, 1)
	<-close
}
