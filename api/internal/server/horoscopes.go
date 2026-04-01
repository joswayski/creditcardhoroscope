package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CreateHoroscopeRequest struct {
	PaymentIntentId string `json:payment_intent_id`
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

	// Check stripe to make sure it's good

}
