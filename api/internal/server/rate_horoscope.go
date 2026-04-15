package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"
)

type RateHoroscopeBody struct {
	PaymentIntentId string `json:"payment_intent_id"`
	Rating          string `json:"rating"`
}

var validRatings = []string{"positive", "negative", "neutral"}

func (s *Server) RateHoroscope(w http.ResponseWriter, r *http.Request) {
	var reqBody RateHoroscopeBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		slog.Error("Error when parsing the request body", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bad request format",
		})
		return
	}

	externalId := r.PathValue("id")

	if !slices.Contains(validRatings, reqBody.Rating) || !strings.HasPrefix(reqBody.PaymentIntentId, "pi_") || externalId == "" {
		slog.Error("Error when parsing the request body", "error", err, "body", reqBody, "externalId", externalId)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bad request, invalid rating, horoscope, or payment intent",
		})
		return
	}

	dbResult, err := s.DB.Exec(r.Context(), `
	UPDATE generations
	SET rating = $1, updated_at = $2
	WHERE
	rating IS NULL AND
	external_id = $3 AND
	payment_intent_id = (SELECT id from payment_intents WHERE payment_intent_id = $4)`,
		reqBody.Rating,
		time.Now().UTC(),
		externalId,
		reqBody.PaymentIntentId,
	)
	if err != nil {
		slog.Error("Error adding rating", "body", reqBody, "externalId", externalId)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "An error ocurred adding your rating, please try again",
		})
		return
	}

	if dbResult.RowsAffected() == 0 {
		slog.Error("Horoscope not found", "body", reqBody, "externalId", externalId)
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Horoscope not found",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Thanks for your feedback! q:^]",
	})

}
