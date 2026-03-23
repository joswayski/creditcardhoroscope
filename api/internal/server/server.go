package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joswayski/creditcardhoroscope/api/internal/config"
)

type Server struct {
	Config     config.Config
	httpServer *http.Server
}

func New(cfg config.Config) *Server {
	s := &Server{Config: cfg}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleSaulGoodman)
	mux.HandleFunc("/health", s.handleSaulGoodman)

	s.httpServer = &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return s
}

func (s *Server) Run() {
	slog.Info(fmt.Sprintf("Server running on port %s", s.Config.Port))
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Error starting API", "error", err)
		os.Exit(1)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
