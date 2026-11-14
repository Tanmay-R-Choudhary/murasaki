package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Unsupported HTTP Method", http.StatusMethodNotAllowed)
	}

	if _, err := w.Write([]byte("PONG")); err != nil {
		http.Error(w, "Something unexpected occured", http.StatusInternalServerError)
		log.Fatal("pongHandler:", err)
	}
}

func main() {
	server := &http.Server{Addr: ":6379"}

	// Add Handlers
	http.HandleFunc("/ping", pingHandler)
	// ------------

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("Exception from server:", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	_, close := context.WithTimeout(context.Background(), 30*time.Second)
	defer close()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Error occured while gracefully shutting down:", err)
		return
	}

	log.Println("Gracefull shutdown complete")
}
