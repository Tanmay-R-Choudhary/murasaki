package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"

	"golang.org/x/sync/errgroup"
)

var payloads []string = []string{"PING", "SOMETHING ELSE"}

func makeConnectionAndSendPing() error {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		return err
	}
	defer conn.Close()

	conn.Write([]byte(payloads[rand.Intn(len(payloads))]))

	buf := make([]byte, 1024) // Max body size restricted to 1KB for now
	if _, err := conn.Read(buf); err != nil {
		return err
	}

	fmt.Println(string(buf))
	return nil
}

func main() {
	errorGroup := new(errgroup.Group)

	for range 10 {
		errorGroup.Go(makeConnectionAndSendPing)
	}

	if err := errorGroup.Wait(); err != nil {
		log.Fatal("Client goroutines erd'd:", err)
	}
}
