package main

import (
	"net"
	"slices"
	"sync"
	"time"
)

type connHandler struct {
	mu    sync.Mutex
	conns []net.Conn
}

func newConnectionHandler() *connHandler {
	return &connHandler{conns: []net.Conn{}}
}

func (ch *connHandler) addConn(conn net.Conn) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.conns = append(ch.conns, conn)
}

func (ch *connHandler) removeConn(conn net.Conn) {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.conns = slices.DeleteFunc(ch.conns, func(ele net.Conn) bool { return conn == ele })
}

func (ch *connHandler) breakAllActiveReads() {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	for _, conn := range ch.conns {
		conn.SetReadDeadline(time.Now())
	}
}
