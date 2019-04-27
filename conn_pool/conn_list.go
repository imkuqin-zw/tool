package conn_pool

import (
	"net"
	"time"
)

type connList struct {
	count       int
	front, back *connNode
}

type connNode struct {
	c       net.Conn
	t       time.Time
	created time.Time
	prev, next *connNode
}

func (cn *connNode) Conn() net.Conn {
	return cn.c
}

func (l *connList) pushFront(cn *connNode) {
	cn.next = l.front
	cn.prev = nil
	if l.count == 0 {
		l.back = cn
	} else {
		l.front.prev = cn
	}
	l.front = cn
	l.count++
	return
}

func (l *connList) popFront() *connNode {
	pc := l.front
	l.count--
	if l.count == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.next.prev = nil
		l.front = pc.next
	}
	pc.next, pc.prev = nil, nil
	return pc
}

func (l *connList) popBack() *connNode {
	pc := l.back
	l.count--
	if l.count == 0 {
		l.front, l.back = nil, nil
	} else {
		pc.prev.next = nil
		l.back = pc.prev
	}
	pc.next, pc.prev = nil, nil
	return pc
}
