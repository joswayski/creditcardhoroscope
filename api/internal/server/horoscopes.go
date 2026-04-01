package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

type CreateHoroscopeRequest struct {
	PaymentIntentId string `json:"payment_intent_id"`
}

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
	defer tx.Rollback(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Sorry, an error ocurred :( Please try again!",
		})
		return
	}

	rows, err := tx.Query(r.Context(), `
	SELECT * from payment_intents where payment_intent_id = $1 FOR UPDATE
	`, req.PaymentIntentId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Sorry, an error ocurred :( Please try again!",
		})
		return
	}
	paymentIntent, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[PaymentIntent])
	if err != nil {
		// TODO add better logging
		if err == pgx.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "Unfortunately, we could not find your payment",
			})
			return
		}

		if paymentIntent.Status != "pending" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{
				// TODO in the future we want to allow completed to be re-used / new used status
				// TODO env email
				"message": "Unfortunately, you've already redeemed your horoscope! Please email contact@josevalerio.com if you have any issues.",
			})
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Sorry, an error ocurred :( Please try again!",
		})
		return
	}

	// Check stripe to make sure it's good

	// TODO debug
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]PaymentIntent{
		"pi": paymentIntent,
	})

}
