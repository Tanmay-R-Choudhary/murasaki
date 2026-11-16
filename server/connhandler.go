package main

import (
	"net"
	"slices"
	"sync"
	"time"
)

type ConnHandler struct {
	mu    sync.Mutex
	conns []net.Conn
}

func NewConnectionHandler() *ConnHandler {
	return &ConnHandler{conns: []net.Conn{}}
}

func (ch *ConnHandler) AddConn(conn net.Conn) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.conns = append(ch.conns, conn)
}

func (ch *ConnHandler) RemoveConn(conn net.Conn) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.conns = slices.DeleteFunc(ch.conns, func(ele net.Conn) bool { return conn == ele })
}

func (ch *ConnHandler) BreakAllActiveReads() {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for _, conn := range ch.conns {
		conn.SetReadDeadline(time.Now())
	}
}
