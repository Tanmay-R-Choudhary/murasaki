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
	"sync"
	"syscall"

	"github.com/Tanmay-R-Choudhary/murasaki/mwp"
	"golang.org/x/net/netutil"
)

func execute(cmd mwp.Data) mwp.Data {
	arr, ok := cmd.(*mwp.MuArray)
	if !ok || len(arr.Value()) == 0 {
		return mwp.NewMuString("ERR invalid command")
	}

	args := arr.Value()
	name, ok := args[0].(*mwp.MuString)
	if !ok {
		return mwp.NewMuString("ERR invalid command")
	}

	switch strings.ToUpper(name.Value()) {
	case "PING":
		return mwp.NewMuString("PONG")
	case "ECHO":
		if len(args) != 2 {
			return mwp.NewMuString("ERR wrong number of arguments for 'echo' command")
		}
		return args[1]
	default:
		return mwp.NewMuArray([]mwp.Data{mwp.NewMuString("Unrecognized")})
	}
}

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

	buf := make([]byte, 1024) // Max body size restricted to 1KB for now
	for {
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
			writeResp(execute(mwpResp).Serialize(), conn)
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
