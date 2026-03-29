package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joswayski/creditcardhoroscope/api/internal/config"
	"github.com/joswayski/creditcardhoroscope/api/internal/database"
	"github.com/joswayski/creditcardhoroscope/api/internal/server"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	apiConfig := config.LoadConfig()

	dbConn := database.Connect(apiConfig.DatabaseURL)
	defer dbConn.Close()

	s := server.New(apiConfig, dbConn)
	go s.Run()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	slog.Info("Shutting down server")

	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := s.Shutdown(shutdownContext)
	if err != nil {
		slog.Error("Error shutting down server", "error", err)
	}

	slog.Info("Server stopped")

}
