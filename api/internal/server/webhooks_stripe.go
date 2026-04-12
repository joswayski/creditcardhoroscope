package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/stripe/stripe-go/v85"
)

func (s *Server) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bad request: no body",
		})
		return
	}

	stripeSignature := r.Header.Get("Stripe-Signature")
	if stripeSignature == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bad request: missing signature",
		})
		return
	}

	event, err := s.Stripe.ConstructEvent(body, stripeSignature, s.Config.StripeWebhookSecretKey)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Invalid request",
		})
		return
	}

	switch event.Type {
	case "charge.refunded":
		// Try to parse it
		var charge stripe.Charge
		err := json.Unmarshal(event.Data.Raw, &charge)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Invalid request - bad data",
			})
			return
		}

		tx, err := s.DB.Begin(r.Context())
		if err != nil {
			slog.Error("Error handling refund event", "error", err, "event", event)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "An error ocurred starting a transsaction when handling a refund event",
			})
			return
		}
		defer tx.Rollback(r.Context())

		// Get the PI in the DB and update its status
		var piId int
		var piStatus string
		err = tx.QueryRow(r.Context(), `
		SELECT id, status FROM payment_intents
		WHERE payment_intent_id = $1 FOR UPDATE
		`, charge.PaymentIntent.ID).Scan(&piId, &piStatus)
		if errors.Is(pgx.ErrNoRows, err) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("Payment intent with ID %s not found", charge.PaymentIntent.ID),
			})
			return
		}

		if piStatus == "refunded" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("Payment intent with ID %s was already refunded", charge.PaymentIntent.ID),
			})
			return
		}

		// Set to refund
		_, err = tx.Exec(r.Context(), `
		UPDATE payment_intents 
		SET status = 'refunded'
		WHERE id = $1`,
			piId)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("Unable to save refund of payment intent %s due to error %e", charge.PaymentIntent.ID, err),
			})
			return
		}

		err = tx.Commit(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"message": fmt.Sprintf("Unable to commit refund of payment intent %s due to error %e", charge.PaymentIntent.ID, err),
			})
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
