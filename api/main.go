package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", hello)

	server := http.Server{
		Addr:         ":8080", // ! todo env
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Error starting API", "error", err)
		return
	}
}

func hello(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "hello world!",
	})
}
