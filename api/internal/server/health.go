package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) SaulGoodman(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Saul Goodman",
	})
}
