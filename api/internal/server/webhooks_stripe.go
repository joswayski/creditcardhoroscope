package server

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

func (s *Server) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	event, err := s.Stripe.ConstructEvent(body, r.Header["Stripe-Signature"][0], s.Config.StripeWebhookSecretKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request",
		})
		return
	}
	slog.Info("Request", "body", r.Body, "headers", r.Header)
}
