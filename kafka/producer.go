package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"time"
)

type ProducerConfig struct {
	Brokers []string
	Sync    bool
	Retries uint32
}

type errHandle func(*sarama.ProducerError)
type sucHandle func(*sarama.ProducerMessage)

type Producer struct {
	sarama.AsyncProducer
	sarama.SyncProducer
	conf    *ProducerConfig
	errDeal errHandle
	sucDeal sucHandle
}

func NewProducer(conf *ProducerConfig, errDeal errHandle, sucDeal sucHandle) (*Producer, error) {
	if conf.Retries == 0 {
		conf.Retries = 3
	}
	p := &Producer{
		conf:    conf,
		errDeal: errDeal,
		sucDeal: sucDeal,
	}
	if !conf.Sync {
		if err := p.asyncDial(); err != nil {
			return nil, err
		}
	} else {
		if err := p.syncDial(); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *Producer) syncDial() (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	for i := uint32(0); i < p.conf.Retries; i++ {
		if p.SyncProducer, err = sarama.NewSyncProducer(p.conf.Brokers, config); err == nil {
			return
		}
		sarama.Logger.Printf("new sync producer fault times(%d) error(%v)", i, err)
		time.Sleep(time.Second)
	}
	return
}

func (p *Producer) asyncDial() (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal     // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms
	for i := uint32(0); i < p.conf.Retries; i++ {
		if p.AsyncProducer, err = sarama.NewAsyncProducer(p.conf.Brokers, config); err == nil {
			go p.errProcess(p.errDeal)
			go p.successProcess(p.sucDeal)
			break
		}
		sarama.Logger.Printf("new async producer fault times(%d) error(%v)", i, err)
		time.Sleep(time.Second)
	}
	return
}

func (p *Producer) errProcess(deal errHandle) {
	err := p.Errors()
	for {
		e, ok := <-err
		if !ok {
			return
		}
		sarama.Logger.Printf("kafka producer send message(%v) failed error(%v)", e.Msg, e.Err)
		if deal != nil {
			deal(e)
		}
	}
}

func (p *Producer) successProcess(deal sucHandle) {
	suc := p.Successes()
	for {
		msg, ok := <-suc
		if !ok {
			return
		}
		sarama.Logger.Printf("kafka producer send message(%v) sucsess", msg)
		if deal != nil {
			deal(msg)
		}
	}
}

func (p *Producer) Input(c context.Context, msg *sarama.ProducerMessage) (err error) {
	if !p.conf.Sync {
		msg.Metadata = c
		p.AsyncProducer.Input() <- msg
	} else {
		if _, _, err = p.SyncProducer.SendMessage(msg); err != nil {
			sarama.Logger.Print("syncProducer send msg(%v) fault error(%v): ", msg, err)
		}
	}
	return
}

func (p *Producer) Close() (err error) {
	if !p.conf.Sync {
		if p.AsyncProducer != nil {
			return p.AsyncProducer.Close()
		}
	}
	if p.SyncProducer != nil {
		return p.SyncProducer.Close()
	}
	return
}
