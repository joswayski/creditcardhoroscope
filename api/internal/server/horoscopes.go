package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joswayski/creditcardhoroscope/api/internal/horoscopes"
	"github.com/joswayski/creditcardhoroscope/api/internal/types"
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
	dbPaymentIntent, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[types.PaymentIntent])
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
	if !dbPaymentIntent.AllowsGenerations() {
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
	card_last_4 = $5, card_country = $6, card_postal = $7, updated_at = $8
	WHERE id = $9`,
		"paid",
		cardDetails.brand,
		cardDetails.expMonth,
		cardDetails.expYear,
		cardDetails.last4,
		cardDetails.country,
		cardDetails.postalCode,
		time.Now().UTC(),
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
			Model:        "poop",
			Instructions: openai.String(s.Config.AISystemPrompt),
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(horoscopes.FormatUserMessage(&dbPaymentIntent)),
			}})

		if aiErr == nil && aiResponse.Status != "failed" {
			// Saul Goodman
			break
		}
		slog.Error("AI Generation failed", "attempt", i, "error", aiErr, "pi", dbPaymentIntent.PaymentIntentID)

		if i >= maxRetries {
			go func() {
				//  Try to write the failure in the background
				_, err = s.DB.Exec(context.Background(), `
					INSERT INTO generations (payment_intent_id, status, error)
					VALUES ($1, $2, $3)
					`, dbPaymentIntent.ID, "failed", aiErr.Error())
				if err != nil {
					slog.Error("Error inserting failed generation", "pi", dbPaymentIntent.PaymentIntentID, "error", aiErr)
				}
			}()

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
	tx, err = s.DB.Begin(r.Context())
	if err != nil {
		slog.Error("Error starting transaction after horoscope was generated", "pi", dbPaymentIntent.PaymentIntentID, "horoscope", horoscope, "aiResponse", aiResponse)
		w.WriteHeader(http.StatusCreated)
		// This is our issue at this point but we can still give the user a good time
		json.NewEncoder(w).Encode(map[string]any{
			"horoscope": horoscope,
		})
		return
	}
	defer tx.Rollback(r.Context())

	// Lock the PI again
	rows, err = tx.Query(r.Context(), `
	SELECT * FROM payment_intents 
	WHERE id = $1 FOR UPDATE
	`, dbPaymentIntent.ID)

	if err != nil {
		slog.Error("Error retrieving horoscope from DB while adding generation", "pi", dbPaymentIntent.PaymentIntentID, "horoscope", horoscope, "error", err, "aiResponse", aiResponse)
		w.WriteHeader(http.StatusCreated)
		// This is our issue at this point but we can still give the user a good time
		json.NewEncoder(w).Encode(map[string]any{
			"horoscope": horoscope,
		})
		return
	}

	dbPaymentIntent, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[types.PaymentIntent])
	if err != nil {
		slog.Error("Error retrieving horoscope from DB while adding generation 2", "pi", dbPaymentIntent.PaymentIntentID, "horoscope", horoscope, "error", err, "aiResponse", aiResponse)
		w.WriteHeader(http.StatusCreated)
		// This is our issue at this point but we can still give the user a good time
		json.NewEncoder(w).Encode(map[string]any{
			"horoscope": horoscope,
		})
		return
	}

	if !dbPaymentIntent.AllowsGenerations() {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Unfortunately, this payment cannot be redeemed for a horoscope. If you have any questions email %s with this ID: %s", s.Config.SupportEmail, dbPaymentIntent.PaymentIntentID),
		})
		return
	}
	// Update generations row
	_, err = tx.Exec(r.Context(), `
	INSERT INTO generations 
	(payment_intent_id, status, or_gen_id, or_model, or_tokens_used, horoscope)
	VALUES ($1, $2, $3, $4, $5, $6)
	`, dbPaymentIntent.ID, "completed", aiResponse.ID, string(aiResponse.Model), aiResponse.Usage.TotalTokens, aiResponse.OutputText())

	if err != nil {
		// Again, this is our problem at this point
		slog.Error("Error adding horoscope to DB during insert", "pi", dbPaymentIntent.PaymentIntentID, "horoscope", horoscope, "error", err, "aiResponse", aiResponse)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"horoscope": horoscope,
		})
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		// Again, this is our problem at this point
		slog.Error("Error adding horoscope to DB during commit", "pi", dbPaymentIntent.PaymentIntentID, "horoscope", horoscope, "error", err, "aiResponse", aiResponse)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{
			"horoscope": horoscope,
		})
		return
	}

	//
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"horoscope": horoscope,
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
