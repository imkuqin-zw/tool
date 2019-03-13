package codis_cli

import (
	"context"
	"encoding/json"
	"fmt"
	etcdV2 "github.com/coreos/etcd/client"
	"github.com/garyburd/redigo/redis"
	"math/rand"
	"sync"
	"time"
)

type Config struct {
	MaxIdle        int
	MaxActive      int
	IdleTimeOut    time.Duration
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RegisterKey    string
}

type NodeConfig struct {
	Token     string `json:"token"`
	ProxyAddr string `json:"proxy_addr"`
	Pwd       string `json:"pwd"`
	Hostname  string `json:"hostname"`
	ProtoType string `json:"proto_type"`
}

type Client struct {
	mu       sync.RWMutex
	nodePool map[string]*redis.Pool
	nodeArr  []string
	nodeLen  uint64
	config   Config
	etcd     etcdV2.KeysAPI
}

func NewClient(config Config, clientV2 etcdV2.Client) *Client {
	return &Client{
		nodePool: make(map[string]*redis.Pool),
		nodeArr:  make([]string, 0),
		etcd:     etcdV2.NewKeysAPI(clientV2),
		config:   config,
	}
}

func (c *Client) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//c.etcd.Put(ctx, "/test", "212")
	res, err := c.etcd.Get(ctx, c.config.RegisterKey,
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
		fmt.Println(nodeConfig.Pwd)
		c.addRedisPool(nodeConfig)
	}

	go c.watch()
	return nil
}

func (c *Client) GetConn() redis.Conn {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.nodeLen == 0 {
		return nil
	}
	index := rand.Uint64() % c.nodeLen
	return c.nodePool[c.nodeArr[index]].Get()
}

func (c *Client) addRedisPool(node NodeConfig) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.nodePool[node.Token]; ok {
		return
	}
	c.nodeLen++
	c.nodeArr = append(c.nodeArr, node.Token)
	c.nodePool[node.Token] = &redis.Pool{
		MaxIdle:     c.config.MaxIdle,
		MaxActive:   c.config.MaxActive,
		IdleTimeout: c.config.IdleTimeOut,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(node.ProtoType, node.ProxyAddr)
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				println("TestOnBorrow", err)
			}
			return err
		},
	}
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
	watcher := c.etcd.Watcher(c.config.RegisterKey, &etcdV2.WatcherOptions{Recursive: true})
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
