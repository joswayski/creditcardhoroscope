package main

import (
	"context"
	"fmt"
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
	err := run()
	if err != nil {
		slog.Error("Fatal error", "error", err)
		os.Exit(1)
	}
}

func run() error {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	apiConfig := config.LoadConfig()

	pool := database.Connect(apiConfig.DatabaseURL)
	defer pool.Close()

	err := database.RunMigrations(pool)
	if err != nil {
		return fmt.Errorf("Error running migrations %w", err)
	}

	s := server.New(apiConfig, pool)
	go s.Run()

	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM)
	<-quitChannel

	slog.Info("Shutting down server")

	shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.Shutdown(shutdownContext)
	if err != nil {
		return fmt.Errorf("Error shutting down server %w", err)
	}

	slog.Info("Server stopped")
	return nil
}
