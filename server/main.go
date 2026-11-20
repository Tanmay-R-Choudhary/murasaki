package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp"
	"golang.org/x/net/netutil"
)

func writeResp(data string, conn net.Conn) {
	if _, err := conn.Write([]byte(data)); err != nil {
		log.Println("Error occured while writing response: ", err)
	}
}

func handle(conn net.Conn, close func()) {
	defer func() {
		close()
		conn.Close()
	}()

	for {
		buf := make([]byte, 1024) // Max body size restricted to 1KB for now
		n, err := conn.Read(buf)
		if err != nil {
			if !(errors.Is(err, io.EOF) || errors.Is(err, os.ErrDeadlineExceeded)) {
				log.Println("Error occured while reading body of request: ", err)
			}
			break
		}

		if mwpResp, _, err := mwp.Deserialize(string(buf[:n])); err != nil {
			log.Println("Error occured while deserializing client command: ", err)
			break
		} else {
			if mwpResp.Equals(mwp.NewMuArray([]mwp.Data{mwp.NewMuString("PING")})) {
				writeResp(mwp.NewMuArray([]mwp.Data{mwp.NewMuString("PONG")}).Serialize(), conn)
			} else {
				writeResp(mwp.NewMuArray([]mwp.Data{mwp.NewMuString("Unrecognized")}).Serialize(), conn)
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Error occured while making listener for tcp server: ", err)
	}

	limitedListener := netutil.LimitListener(listener, 50)

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigchan
		// Add graceful shutdown handling here
		limitedListener.Close()
	}()

	wg := &sync.WaitGroup{}
	ch := NewConnectionHandler()

	for {
		conn, err := limitedListener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				log.Println("Closing server, not accepting anymore connections")
				break
			}

			log.Println("Error occured while accepting connection: ", err)
			continue
		}

		ch.AddConn(conn)
		wg.Go(func() {
			handle(conn, func() {
				ch.RemoveConn(conn)
			})
		})
	}

	ch.BreakAllActiveReads()
	wg.Wait()
	fmt.Println("Graceful shutdown complete")
}
