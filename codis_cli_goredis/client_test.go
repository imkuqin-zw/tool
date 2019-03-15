package codis_cli_redigo

import (
	"fmt"
	etcdV2 "github.com/coreos/etcd/client"
	"github.com/go-redis/redis"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"
)

var client *Client
var etcdClient etcdV2.Client

func init() {
	var err error
	config := Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		RegisterKey:  "/codis3/codis-demo/proxy/",
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
	log.Println("set success")
}

func TestClient_Set(t *testing.T) {
	client := client.GetClient()
	if client == nil {
		fmt.Println("not client")
		return
	}
	defer client.Close()
	err := client.Set("testfdsaf_rertr", "123", time.Duration(0)).Err()
	if err != nil {
		log.Print("error", err.Error())
		return
	}
}
func TestClient_Get(t *testing.T) {
	for i := 0; i < 10; i++ {
		go func() {
			testqps()
		}()
	}
	close := make(chan bool, 1)
	<-close
}

func testqps() {
	for {
		client := client.GetClient()
		if client == nil {
			fmt.Println("not client")
			return
		}
		pipe := client.Pipeline()
		pipe.Process(client.Get("testfdsaf_rertr"))
		pipe.Process(client.Set("testfdsaf_rertr", rand.Int(), time.Duration(0)))
		results, err := pipe.Exec()
		if err != nil {
			log.Fatal("pipe Exec fault", err.Error())
		}
		res1, err := results[0].(*redis.StringCmd).Result()
		if err != nil {
			log.Print("error", err.Error())
			return
		}

		res2, err := results[1].(*redis.StatusCmd).Result()
		if err != nil {
			log.Print("error", err.Error())
			return
		}
		fmt.Println(res1, res2)

		//time.Sleep(100 * time.Millisecond)
	}
}
