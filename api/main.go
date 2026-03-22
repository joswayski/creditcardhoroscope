package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)

	server := http.Server{
		Addr:         ":8080", // todo env
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go startServer(&server)

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	slog.Info("Shutting down server")

	shutdownContext, shutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdown()

	err := server.Shutdown(shutdownContext)
	if err != nil {
		slog.Error("Error shutting down server", "error", err)
	}

	slog.Info("Server stopped")

}

func startServer(server *http.Server) {
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Error starting API", "error", err)
		os.Exit(1)
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "hello world!",
	})
}
