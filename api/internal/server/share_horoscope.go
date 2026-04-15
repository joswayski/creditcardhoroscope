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
	IsPublic        bool   `json:"is_public"`
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

	if !strings.HasPrefix(reqBody.PaymentIntentId, "pi_") || externalId == "" || reqBody.IsPublic != true {
		slog.Error("Error when parsing the request body", "error", err, "body", reqBody, "externalId", externalId)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Bad request, invalid is public, horoscope, or payment intent",
		})
		return
	}

	dbResult, err := s.DB.Exec(r.Context(), `
	UPDATE generations
	SET is_public = $1, updated_at = $2
	WHERE
	is_public IS false AND
	external_id = $3 AND
	payment_intent_id = (SELECT id from payment_intents WHERE payment_intent_id = $4)`,
		reqBody.IsPublic,
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message":        "TODO send link here",
		"shareable_link": fmt.Sprintf("https://TODOADDWEBURLHERE.com/share/%s", externalId),
	})

}
