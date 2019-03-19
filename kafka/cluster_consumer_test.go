package kafka

import (
	"fmt"
	"os"
	"sync"
	"testing"
)

var clusterConfig *ClusterConfig

func init() {
	clusterConfig = &ClusterConfig{
		GroupId: "test3",
		Brokers: []string{"192.168.2.118:9091", "192.168.2.118:9092", "192.168.2.118:9093"},
		Topics:  []string{"test"},
		Offset:  true,
	}
}

func TestClusterConsumer(t *testing.T) {
	var wg = &sync.WaitGroup{}
	wg.Add(2)
	go runClusterConsumer("c_1")
	go runClusterConsumer("c_2")
	wg.Wait()
}

func runClusterConsumer(flag string) {
	c, err := NewClusterConsumer(clusterConfig, nil, nil)
	if err != nil {
		fmt.Println("NewClusterConsumer", err.Error())
		return
	}
	defer c.Close()
	for {
		select {
		case msg, ok := <-c.Messages():
			if ok {
				fmt.Fprintf(os.Stdout, "[%s] %s:%s/%d/%d\t%s\t%s\n", flag, clusterConfig.GroupId, msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
				c.MarkOffset(msg, "") // mark message as processed
				continue
			}
			return
		}
	}
}
