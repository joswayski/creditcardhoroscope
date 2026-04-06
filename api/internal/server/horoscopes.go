package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
	"github.com/stripe/stripe-go/v85"
)

type CreateHoroscopeRequest struct {
	PaymentIntentId string `json:"payment_intent_id"`
}

const maxRetries = 3

func (s *Server) CreateHoroscope(w http.ResponseWriter, r *http.Request) {
	// Check we received something
	var req CreateHoroscopeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || !strings.HasPrefix(req.PaymentIntentId, "pi_") {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Invalid payment_intent_id in the request. Received: '%s'", req.PaymentIntentId),
		})
		return
	}

	// Check our DB to make sure it's still pending
	// TODO in the future we will allow multiple uses
	tx, err := s.DB.Begin(r.Context())
	if err != nil {
		slog.Error("Error starting transaction", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Sorry, an error ocurred :( Please try again!",
		})
		return
	}
	defer tx.Rollback(r.Context())

	rows, err := tx.Query(r.Context(), `
	SELECT * from payment_intents where payment_intent_id = $1 FOR UPDATE
	`, req.PaymentIntentId)
	if err != nil {
		slog.Error("Error querying", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Sorry, an error ocurred :( Please try again!",
		})
		return
	}
	dbPaymentIntent, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[PaymentIntent])
	if err != nil {
		slog.Error("Error getting payment intent", "error", err)
		// TODO add better logging
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Unfortunately, we could not find your payment",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Sorry, an error ocurred :( Please try again!",
		})
		return
	}

	// Check the DB status before proceeding
	// Pending gets let through because we're awaiting a generation
	// Paid gets let through because we'll allow multiple generations (TODO)
	if dbPaymentIntent.Status != "pending" && dbPaymentIntent.Status != "paid" {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Unfortunately, this payment cannot be redeemed for a horoscope. If you have any questions email %s with this ID: %s", s.Config.SupportEmail, dbPaymentIntent.PaymentIntentID),
		})
		return
	}

	stripePaymentIntent, err := s.Stripe.V1PaymentIntents.Retrieve(r.Context(), dbPaymentIntent.PaymentIntentID, &stripe.PaymentIntentRetrieveParams{
		Expand: stripe.StringSlice([]string{"payment_method"}),
	})

	if err != nil || stripePaymentIntent == nil || stripePaymentIntent.PaymentMethod == nil || stripePaymentIntent.PaymentMethod.Card == nil {
		slog.Error("Error retrieving payment intent / method / card", "error", err, "pi", stripePaymentIntent)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Sorry, an error ocurred :( Please try again!",
		})
		return
	}

	if stripePaymentIntent.Status != stripe.PaymentIntentStatusSucceeded {
		// TODO handle others
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Payment has not been received yet! Please try again later",
		})
		return
	}

	cardDetails := getCardDetails(stripePaymentIntent.PaymentMethod)

	// If it is paid, we should update the DB to paid so that the user can retry if the AI fails
	_, err = tx.Exec(r.Context(), `
	UPDATE payment_intents
	SET status = $1, card_brand = $2, card_exp_month = $3, card_exp_year = $4,
	card_last_4 = $5, card_country = $6, card_postal = $7
	WHERE id = $8`,
		"paid",
		cardDetails.brand,
		cardDetails.expMonth,
		cardDetails.expYear,
		cardDetails.last4,
		cardDetails.country,
		cardDetails.postalCode,
		dbPaymentIntent.ID)

	if err != nil {
		slog.Error("Error saving horoscope paid status to DB", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "An error ocurred creating your horoscope, please try again",
		})
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		slog.Error("Error saving horoscope paid status to DB", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "An error ocurred creating your horoscope, please try again",
		})
		return
	}

	var aiResponse *responses.Response
	var aiErr error
	for i := 1; i < maxRetries+1; i++ {
		// Call AI API
		// It is recommended to use a fast model here
		// To not give the presense of the streaming/ai interface vibes
		// Should look more like a finished response
		aiResponse, aiErr = s.AI.Responses.New(r.Context(), responses.ResponseNewParams{
			Model:        s.Config.AIModel,
			Instructions: openai.String(s.Config.AISystemPrompt),
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String("user details"),
			}})

		if err == nil && aiResponse.Status != "failed" {
			// Saul Goodman
			break
		}
		slog.Error("AI Generation failed", "attempt", i, "error", aiErr, "pi", dbPaymentIntent.PaymentIntentID)

		if i >= maxRetries {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{
				"message": `Unfortunately, we could not generate a horoscope for you.`,
			})
			// TODO refund
			return
		}

		// Retry
		delay := 200 * math.Pow(2, float64(i))
		time.Sleep(time.Millisecond * time.Duration(delay))
	}

	horoscope := aiResponse.OutputText()

	// TODO debug
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"db_pi":     dbPaymentIntent,
		"stripe_pi": stripePaymentIntent,
	})

}

type cardInfo struct {
	brand      string
	expMonth   string
	expYear    string
	last4      string
	country    string
	postalCode *string
}

func getCardDetails(pm *stripe.PaymentMethod) cardInfo {
	var postalCode *string
	if pm.BillingDetails != nil && pm.BillingDetails.Address != nil && pm.BillingDetails.Address.PostalCode != "" {
		pc := pm.BillingDetails.Address.PostalCode
		postalCode = &pc
	}

	return cardInfo{
		brand:      string(pm.Card.Brand),
		expMonth:   fmt.Sprintf("%d", pm.Card.ExpMonth),
		expYear:    fmt.Sprintf("%d", pm.Card.ExpYear),
		last4:      pm.Card.Last4,
		country:    pm.Card.Country,
		postalCode: postalCode,
	}
}
