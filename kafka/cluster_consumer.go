package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"go.uber.org/zap"
	"time"
)

type clusterErrHandle func(error)
type clusterNtyHandle func(notification *cluster.Notification)

type ClusterConfig struct {
	Brokers  []string
	Topics   []string
	GroupId  string
	Offset   bool
	Retries  uint32
	LogDebug bool
}

type ClusterConsumer struct {
	conf *ClusterConfig
	*cluster.Consumer
	errDeal    clusterErrHandle
	notifyDeal clusterNtyHandle
}

func NewClusterConsumer(config *ClusterConfig, notifyDeal clusterNtyHandle, errDeal clusterErrHandle) (*ClusterConsumer, error) {
	if config.Retries == 0 {
		config.Retries = 3
	}
	consumer := &ClusterConsumer{
		conf:       config,
		errDeal:    errDeal,
		notifyDeal: notifyDeal,
	}
	if err := consumer.dial(); err != nil {
		return nil, err
	}
	return consumer, nil
}

func (c *ClusterConsumer) dial() (err error) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	if c.conf.Offset {
		config.Consumer.Offsets.Initial = sarama.OffsetNewest
	} else {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}
	for i := uint32(0); i < c.conf.Retries; i++ {
		c.Consumer, err = cluster.NewConsumer(c.conf.Brokers, c.conf.GroupId, c.conf.Topics, config)
		if err == nil {
			go c.notificationProcess(c.notifyDeal)
			go c.errProcess(c.errDeal)
			return
		}

		Log.Warn("NewSyncProducer fault", zap.Uint32("time", i), zap.Error(err))
		time.Sleep(time.Second)
	}
	return
}

func (c *ClusterConsumer) notificationProcess(deal clusterNtyHandle) {
	notify := c.Notifications()
	for {
		ntf, ok := <-notify
		if !ok {
			return
		}
		Log.Debug("kafka cluster consumer notification send", zap.Any("notify", ntf))
		if deal != nil {
			deal(ntf)
		}
	}
}

func (c *ClusterConsumer) errProcess(deal clusterErrHandle) {
	err := c.Errors()
	for {
		e, ok := <-err
		if !ok {
			return
		}
		Log.Error("kafka cluster consumer", zap.Any("group_id", c.conf.GroupId), zap.Error(e))
		if deal != nil {
			deal(e)
		}
	}
}

func (c *ClusterConsumer) Close() error {
	if c.Consumer != nil {
		return c.Consumer.Close()
	}
	return nil
}
