package server

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type ShareHoroscopeBody struct {
	PaymentIntentId string `json:"payment_intent_id"`
}

func (s *Server) ShareHoroscope(w http.ResponseWriter, r *http.Request) {
	var reqBody ShareHoroscopeBody
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

	if !strings.HasPrefix(reqBody.PaymentIntentId, "pi_") || externalId == "" {
		slog.Error("Error when parsing the request body", "error", err, "body", reqBody, "externalId", externalId)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bad request: invalid  horoscope, or payment intent",
		})
		return
	}

	dbResult, err := s.DB.Exec(r.Context(), `
	UPDATE generations
	SET is_public = true, updated_at = $1
	WHERE
	external_id = $2 AND
	payment_intent_id = (SELECT id from payment_intents WHERE payment_intent_id = $3)`,
		time.Now().UTC(),
		externalId,
		reqBody.PaymentIntentId,
	)
	if err != nil {
		slog.Error("Error changing visibility", "body", reqBody, "externalId", externalId)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "An error ocurred sharing your horoscope, please try again",
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

	shareableLink := fmt.Sprintf("%s/%s", s.Config.BaseURL, externalId)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":        fmt.Sprintf("Your horoscope is now public! View it here: %s", shareableLink),
		"shareable_link": shareableLink,
	})

}
