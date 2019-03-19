package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"strconv"
	"testing"
	"time"
)

var producerConfig *ProducerConfig

func init() {
	producerConfig = &ProducerConfig{
		Brokers: []string{"192.168.2.118:9091", "192.168.2.118:9092", "192.168.2.118:9093"},
		Sync:    false,
	}
}

func TestProducer(t *testing.T) {
	producer, err := NewProducer(producerConfig, nil, nil)
	if err != nil {
		fmt.Println("NewProducer", err.Error())
		return
	}
	var value []byte
	var i int
	for {
		value = []byte(strconv.Itoa(i))
		msg := &sarama.ProducerMessage{
			Topic: "test",
			Key:   sarama.StringEncoder(fmt.Sprintf("single_msg_%d", i)),
			Value: sarama.ByteEncoder(value),
		}
		i++
		producer.Input(context.Background(), msg)
		time.Sleep(time.Millisecond * 100)
	}
	producer.Close()
}
