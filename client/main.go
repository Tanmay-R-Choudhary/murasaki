package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
)

var payloads []string = []string{"PING", "SOMETHING ELSE"}

func makeStandaloneConnectionAndSendPing() error {
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

func makeManyRequestsWithOneConnection() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error occured while opening tcp connection: ", err)
	}
	defer conn.Close()

	for range 10 {
		conn.Write([]byte("PING"))

		buf := make([]byte, 1024) // Max body size restricted to 1KB for now
		if _, err := conn.Read(buf); err != nil {
			log.Fatal("Error occured while reading response from server: ", err)
			break
		}

		fmt.Println(string(buf))
	}
}

func main() {
	makeManyRequestsWithOneConnection()
}
