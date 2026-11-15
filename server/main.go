package main

import (
	"errors"
	"io"
	"log"
	"net"

	"golang.org/x/net/netutil"
)

func writeResp(data string, conn net.Conn) {
	if _, err := conn.Write([]byte(data)); err != nil {
		log.Println("Error occured while writing response: ", err)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024) // Max body size restricted to 1KB for now
		n, err := conn.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("Closed connection")
				break
			} else {
				log.Println("Error occured while reading body of request: ", err)
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
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error occured while making listener for tcp server: ", err)
	}

	limitedListener := netutil.LimitListener(listener, 50)

	for {
		conn, err := limitedListener.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				log.Println("Error occured while accepting connection: ", err)
				break
			}
			continue
		}

		go handle(conn)
	}

	log.Println("Graceful shutdown complete")
}
