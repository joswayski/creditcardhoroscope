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
	"github.com/joswayski/creditcardhoroscope/api/internal/webhooks"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/stripe/stripe-go/v85"
)

type Server struct {
	Config     config.Config
	httpServer *http.Server
	DB         *pgxpool.Pool
	Stripe     *stripe.Client
	cancel     context.CancelFunc
	AI         openai.Client
}

func New(cfg config.Config, pool *pgxpool.Pool) *Server {
	s := &Server{Config: cfg, DB: pool, Stripe: stripe.NewClient(cfg.StripeSecretKey), AI: openai.NewClient(
		option.WithAPIKey(cfg.AIAPIKey),
		option.WithBaseURL(cfg.AIBaseURL),
	)}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", s.Root)
	mux.HandleFunc("GET /api/v1/{$}", s.Root)
	mux.HandleFunc("GET /api/v1/health", s.SaulGoodman)

	piRateLimiter := middleware.CreateRateLimiter(time.Second*5, 2)
	go piRateLimiter.BackgroundCleanup(ctx)
	mux.HandleFunc("POST /api/v1/payment-intents", middleware.BodySize(middleware.RateLimit(piRateLimiter, s.CreatePaymentIntent), 0))

	// TODO In the future we will allow multiple generations. For now to stop some spam
	horoscopeRateLimiter := middleware.CreateRateLimiter(time.Second*5, 2)
	go horoscopeRateLimiter.BackgroundCleanup(ctx)
	mux.HandleFunc("POST /api/v1/horoscopes", middleware.BodySize(middleware.RateLimit(horoscopeRateLimiter, s.CreateHoroscope), 512))

	mux.HandleFunc("POST /api/v1/webhooks/stripe", middleware.BodySize(middleware.IPWhitelist(s.StripeWebhook, webhooks.StripeIps), 69420))

	// Catchall
	mux.HandleFunc("/", s.FourOhFour)

	s.httpServer = &http.Server{
		Addr: ":" + cfg.Port,
		// left to right
		Handler:      middleware.CORS(middleware.JSONResponseHeader(mux)),
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
	s.cancel()
	return s.httpServer.Shutdown(ctx)
}
