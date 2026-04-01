package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/stripe/stripe-go/v85"
)

type PaymentIntent struct {
	ID              int64     `db:"id"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
	PaymentIntentID string    `db:"payment_intent_id"`
	Amount          int       `db:"amount"`
	Currency        string    `db:"currency"`
	Status          string    `db:"status"`
	CardBrand       *string   `db:"card_brand"`
	CardExpMonth    *string   `db:"card_exp_month"`
	CardExpYear     *string   `db:"card_exp_year"`
	CardLast4       *string   `db:"card_last_4"`
	CardCountry     *string   `db:"card_country"`
	CardPostal      *string   `db:"card_postal"`
}

func (s *Server) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	// Create the PI in stripe
	paymentIntentParams := &stripe.PaymentIntentCreateParams{
		Amount:   stripe.Int64(100), // Always $1
		Currency: stripe.String(stripe.CurrencyUSD),
		AutomaticPaymentMethods: &stripe.PaymentIntentCreateAutomaticPaymentMethodsParams{
			Enabled: new(true),
		},
	}

	piResult, err := s.Stripe.V1PaymentIntents.Create(r.Context(), paymentIntentParams)
	if err != nil {
		slog.Error("Error creating payment intent", "error", err, "params", paymentIntentParams)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "An error ocurred loading your payment form, please try again later :(",
			"error":   "69",
		})
		return
	}

	// Insert the PI in the DB
	_, err = s.DB.Exec(r.Context(), `
	INSERT INTO payment_intents (payment_intent_id, amount, currency)
	VALUES ($1, $2, $3)`, piResult.ID, piResult.Amount, piResult.Currency)

	if err != nil {
		// Try to cancel it in stripe in the background
		slog.Error("Error saving payment intent to the DB", "error", err, "params", paymentIntentParams)
		go func() {
			_, err := s.Stripe.V1PaymentIntents.Cancel(context.Background(), piResult.ID, &stripe.PaymentIntentCancelParams{
				CancellationReason: stripe.String("abandoned"),
			})
			if err != nil {
				slog.Error("Error cancelling payment intent after DB write failure", "error", err, "params", paymentIntentParams, "pi", piResult)

			}
		}()

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "An error ocurred loading your payment form, please try again later :(",
			"error":   "67"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message":           "Payment intent created!",
		"client_secret":     piResult.ClientSecret,
		"payment_intent_id": piResult.ID,
	})
}
