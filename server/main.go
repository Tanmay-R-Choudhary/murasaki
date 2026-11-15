package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/net/netutil"
)

func writeResp(data string, conn net.Conn) {
	if _, err := conn.Write([]byte(data)); err != nil {
		log.Fatal("Error occured while writing response: ", err)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024) // Max body size restricted to 1KB for now
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Print("Closed connection")
				break
			} else {
				log.Fatal("Error occured while reading body of request: ", err)
				return
			}
		}

		if string(buf[:n]) == "PING" {
			writeResp("+PONG\r\n", conn)
		} else {
			writeResp("Unrecognized", conn)
		}
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
