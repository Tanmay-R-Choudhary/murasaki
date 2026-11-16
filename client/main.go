package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func sendCommand(cmd string, conn net.Conn) error {
	if _, err := conn.Write([]byte(cmd)); err != nil {
		log.Println("Error occured while writing: ", err)
		return err
	}

	buf := make([]byte, 1024) // Max body size restricted to 1KB for now
	if _, err := conn.Read(buf); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		} else {
			return err
		}
	}

	fmt.Println(string(buf))
	return nil
}

func main() {
	if len(os.Args) < 1 {
		log.Fatal("Expected input command")
	}

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println("Error occured while opening connection to server: ", err)
		return
	}
	defer conn.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		var rawCommand string
		for {
			_, err := fmt.Scan(&rawCommand)
			if err != nil {
				if errors.Is(err, io.EOF) {
					sigchan <- syscall.SIGTERM
					break
				} else {
					log.Println("Error occured while reading input: ", err)
					continue
				}
			}
			for v := range strings.SplitSeq(rawCommand, "\\n") {
				sendCommand(strings.ToUpper(v), conn)
			}
		}
	}()

	<-sigchan // Waiting for interrupt signal from os
	log.Print("Client shutdown triggered")
}
