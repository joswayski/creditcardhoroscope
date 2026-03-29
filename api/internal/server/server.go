package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joswayski/creditcardhoroscope/api/internal/config"
	"github.com/joswayski/creditcardhoroscope/api/internal/middleware"
	"github.com/stripe/stripe-go/v85"
)

type Server struct {
	Config     config.Config
	httpServer *http.Server
	DB         *pgxpool.Pool
	Stripe     *stripe.Client
}

func New(cfg config.Config, pool *pgxpool.Pool) *Server {
	s := &Server{Config: cfg, DB: pool, Stripe: stripe.NewClient(cfg.StripeSecretKey)}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.Root)
	mux.HandleFunc("GET /api/v1/{$}", s.Root)
	mux.HandleFunc("GET /api/v1/health", s.SaulGoodman)
	mux.HandleFunc("POST /api/v1/payment-intents", s.CreatePaymentIntent)
	mux.HandleFunc("POST /api/v1/horoscopes", s.CreateHoroscope)

	// Catchall
	mux.HandleFunc("/", s.FourOhFour)

	s.httpServer = &http.Server{
		Addr: ":" + cfg.Port,
		// left to right
		Handler:      middleware.CORS(middleware.JSONHeader(mux)),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	return s
}

func (s *Server) Run() {
	slog.Info(fmt.Sprintf("Server running on http://localhost:%s", s.Config.Port))
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Error starting API", "error", err)
		os.Exit(1)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
