package fastdfs

import (
	"net"
	"time"
)

type connList struct {
	count       int
	front, back *ConnNode
}

type ConnNode struct {
	c       net.Conn
	t       time.Time
	created time.Time
	prev, next *ConnNode
}

func (cn *ConnNode) Conn() net.Conn {
	return cn.c
}

func (l *connList) pushFront(cn *ConnNode) {
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

func (l *connList) popFront() *ConnNode {
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

func (l *connList) popBack() *ConnNode {
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
