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

	"github.com/Tanmay-R-Choudhary/murasaki/mwp"
)

func strToMwpRepr(cmd string) mwp.Data {
	cmdElem := []mwp.Data{}

	// Assuming for now command contains only strings
	for v := range strings.SplitSeq(cmd, " ") {
		cmdElem = append(cmdElem, mwp.NewMuString(v))
	}

	return mwp.NewMuArray(cmdElem)
}

func sendCommand(cmd string, conn net.Conn) error {
	mwpRepr := strToMwpRepr(cmd)
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
			recursiveDebugPrint(mwpResp)
		}
	}

	return nil
}

func recursiveDebugPrint(s mwp.Data) {
	switch s.Type() {
	case mwp.MuStringType:
		if v, ok := s.(*mwp.MuString); ok {
			fmt.Println(v.Value())
		} else {
			fmt.Println("ERROR")
		}
	case mwp.MuIntType:
		if v, ok := s.(*mwp.MuInt); ok {
			fmt.Println(v.Value())
		} else {
			fmt.Println("ERROR")
		}
	case mwp.MuArrayType:
		if v, ok := s.(*mwp.MuArray); ok {
			for _, d := range v.Value() {
				recursiveDebugPrint(d)
			}
		} else {
			fmt.Println("ERROR")
		}
	}
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
			for v := range strings.SplitSeq(rawCommand, "\n") {
				sendCommand(strings.ToUpper(v), conn)
			}
		}
	}()

	<-sigchan // Waiting for interrupt signal from os
	log.Print("Client shutdown triggered")
}
