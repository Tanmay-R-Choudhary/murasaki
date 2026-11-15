package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/netutil"
)

func handle(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024) // Max body size restricted to 1KB for now
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatal("Error occured while reading body of request: ", err)
		return
	}

	if string(buf[:n]) == "PING" {
		conn.Write([]byte("PONG"))
		return
	} else {
		conn.Write([]byte("Unrecognized"))
	}
}

func main() {
	serverState := NewServerState()

	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error occured while making listener for tcp server: ", err)
		return
	}

	limitedListener := netutil.LimitListener(listener, 1)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigchan
		limitedListener.Close()
		// Add graceful shutdown handling here
		serverState.setShutdown(true)
	}()

	for {
		// Check if sever is shutting down
		if serverState.getShutdown() {
			break
		}

		conn, err := limitedListener.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Fatal("Error occured while accepting connection: ", err)
			}
			continue
		}

		go handle(conn)
	}

	fmt.Println("Graceful shutdown complete")
}
