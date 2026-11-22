package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp"
)

func sendCommand(cmd string, conn net.Conn) error {
	mwpRepr, err := strToMwpRepr(cmd)
	if err != nil {
		return err
	}

	transferBytes := mwpRepr.Serialize()

	if _, err := conn.Write([]byte(transferBytes)); err != nil {
		log.Println("Error occured while writing: ", err)
		return err
	}

	buf := make([]byte, 1024) // Max body size restricted to 1KB for now
	if n, err := conn.Read(buf); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		} else {
			return err
		}
	} else {
		if mwpResp, _, err := mwp.Deserialize(string(buf[:n])); err != nil {
			return err
		} else {
			fmt.Println(mwpResp)
		}
	}

	return nil
}

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Println("Error occured while opening connection to server: ", err)
		return
	}
	defer conn.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			if err := sendCommand(line, conn); err != nil {
				log.Println("Error occured while sending command: ", err)
				continue
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println("Error occured while reading input: ", err)
		}

		sigchan <- syscall.SIGTERM
	}()

	<-sigchan // Waiting for interrupt signal from os
	log.Print("Client shutdown triggered")
}
