package codis_cli_redigo

import (
	"context"
	"encoding/json"
	etcdV2 "github.com/coreos/etcd/client"
	"github.com/go-redis/redis"
	"math/rand"
	"sync"
	"time"
)

type Config struct {
	DB           int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	PoolSize     int
	PoolTimeout  time.Duration
	RegisterKey  string
}

type NodeConfig struct {
	Token     string `json:"token"`
	ProxyAddr string `json:"proxy_addr"`
	Pwd       string `json:"pwd"`
	Hostname  string `json:"hostname"`
	ProtoType string `json:"proto_type"`
}

type Client struct {
	mu          sync.RWMutex
	nodePool    map[string]*redis.Client
	nodeArr     []string
	nodeLen     uint64
	options     redis.Options
	registerKey string
	etcd        etcdV2.KeysAPI
}

func NewClient(config Config, clientV2 etcdV2.Client) *Client {
	client := &Client{
		nodePool:    make(map[string]*redis.Client),
		nodeArr:     make([]string, 0),
		etcd:        etcdV2.NewKeysAPI(clientV2),
		registerKey: config.RegisterKey,
		options: redis.Options{
			DB: config.DB,
		},
	}
	if config.DialTimeout != 0 {
		client.options.DialTimeout = config.DialTimeout
	}
	if config.ReadTimeout != 0 {
		client.options.ReadTimeout = config.ReadTimeout
	}
	if config.WriteTimeout != 0 {
		client.options.WriteTimeout = config.WriteTimeout
	}
	if config.IdleTimeout != 0 {
		client.options.IdleTimeout = config.IdleTimeout
	}
	if config.PoolSize != 0 {
		client.options.PoolSize = config.PoolSize
	}
	if config.PoolTimeout != 0 {
		client.options.PoolTimeout = config.PoolTimeout
	}
	return client
}

func (c *Client) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := c.etcd.Get(ctx, c.registerKey,
		&etcdV2.GetOptions{Recursive: true},
	)
	if err != nil {
		return err
	}
	for _, node := range res.Node.Nodes {
		nodeConfig := NodeConfig{}
		if err := json.Unmarshal([]byte(node.Value), &nodeConfig); err != nil {
			continue
		}
		c.addRedisPool(nodeConfig)
	}

	go c.watch()
	return nil
}

func (c *Client) GetClient() *redis.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.nodeLen == 0 {
		return nil
	}
	index := rand.Uint64() % c.nodeLen
	return c.nodePool[c.nodeArr[index]]
}

func (c *Client) addRedisPool(node NodeConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.nodePool[node.Token]; ok {
		return
	}
	c.nodeLen++
	c.nodeArr = append(c.nodeArr, node.Token)
	var options redis.Options
	options = c.options
	options.Network = node.ProtoType
	options.Addr = node.ProxyAddr
	c.nodePool[node.Token] = redis.NewClient(&options)
}

func (c *Client) removeRedisPool(node NodeConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.nodePool[node.Token]
	if !ok {
		return
	}
	c.nodeLen--
	c.nodePool[node.Token].Close()
	c.nodePool[node.Token] = nil
	delete(c.nodePool, node.Token)
	for i, item := range c.nodeArr {
		if node.Token == item {
			c.nodeArr = append(c.nodeArr[:i], c.nodeArr[i:]...)
		}
	}
}

func (c *Client) watch() {
	// watch key 监听节点
	watcher := c.etcd.Watcher(c.registerKey, &etcdV2.WatcherOptions{Recursive: true})
	for {
		resp, err := watcher.Next(context.Background())
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		nodeConfig := NodeConfig{}
		if err := json.Unmarshal([]byte(resp.Node.Value), &nodeConfig); err != nil {
			continue
		}
		switch resp.Action {
		case "set":
			c.addRedisPool(nodeConfig)
		case "delete":
			c.removeRedisPool(nodeConfig)
		}
	}
}
