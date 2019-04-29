package fastdfs

import (
	"errors"
	"sync"
	"time"
)

var ErrPoolExhausted = errors.New("connection conn_pool exhausted")

type DialFunc func() (*conn, error)

type TestFunc func(c *conn, t time.Time) error

type Pool interface {
	Get() (*ConnNode, error)

	Put(*ConnNode, bool) error

	Close()

	ActiveNum() int
}

type pool struct {
	dial DialFunc

	testOnBorrow TestFunc

	maxIdle int

	maxActive int

	active int

	idleTimeout time.Duration

	maxConnLifetime time.Duration

	mu sync.Mutex

	closed bool

	ch chan struct{}

	idle connList
}

type Config struct {
	MaxIdle int

	MaxActive int

	IdleTimeout time.Duration
}

func NewPool(conf Config, dial DialFunc, test TestFunc) Pool {
	pool := &pool{
		dial:         dial,
		testOnBorrow: test,
		maxIdle:      conf.MaxIdle,
		maxActive:    conf.MaxActive,
		idleTimeout:  conf.IdleTimeout,
	}
	pool.init()
	return pool
}

func (p *pool) init() {
	p.ch = make(chan struct{}, p.maxActive)
	for i := 0; i < p.maxActive; i++ {
		p.ch <- struct{}{}
	}
}

func (p *pool) Get() (*ConnNode, error) {
	<-p.ch
	p.mu.Lock()
	if p.idleTimeout > 0 {
		n := p.idle.count
		for i := 0; i < n && p.idle.back != nil && p.idle.back.t.Add(p.idleTimeout).Before(time.Now()); i++ {
			cn := p.idle.popBack()
			p.mu.Unlock()
			cn.c.Close()
			p.mu.Lock()
			p.active--
		}
	}

	for {
		cn := p.idle.popFront()
		if cn == nil {
			break
		}
		p.mu.Unlock()
		if (p.testOnBorrow == nil || p.testOnBorrow(cn.c, cn.t) == nil) &&
			(p.maxConnLifetime == 0 || time.Now().Sub(cn.created) < p.maxConnLifetime) {
			return cn, nil
		}
		cn.c.Close()
		p.mu.Lock()
		p.active--
	}

	if p.closed {
		p.mu.Unlock()
		return nil, errors.New("get on closed conn_pool")
	}

	if p.active >= p.maxActive {
		p.mu.Unlock()
		return nil, ErrPoolExhausted
	}

	p.active++
	c, err := p.dial()
	if err != nil {
		c = nil
		p.mu.Lock()
		p.active--
		if p.ch != nil && !p.closed {
			p.ch <- struct{}{}
		}
		p.mu.Unlock()
	}
	return &ConnNode{c: c, created: time.Now()}, err
}

func (p *pool) Put(cn *ConnNode, forceClose bool) error {
	p.mu.Lock()
	if !p.closed && !forceClose {
		cn.t = time.Now()
		if p.idle.count < p.maxIdle {
			p.idle.pushFront(cn)
			cn = nil
		}
	}

	if cn != nil {
		p.mu.Unlock()
		cn.c.Close()
		p.mu.Lock()
		p.active--
	}

	if p.ch != nil && !p.closed {
		p.ch <- struct{}{}
	}
	p.mu.Unlock()
	return nil
}

func (p *pool) ActiveNum() int {
	return p.active
}

func (p *pool) Close() {
	p.mu.Lock()
	if p.closed {
		p.mu.Unlock()
		return
	}
	p.closed = true
	p.active -= p.idle.count
	headCN := p.idle.popFront()
	p.idle.count = 0
	p.idle.front, p.idle.back = nil, nil
	if p.ch != nil {
		close(p.ch)
	}
	p.mu.Unlock()
	for ; headCN != nil; headCN = headCN.next {
		headCN.c.Close()
	}
	return
}
