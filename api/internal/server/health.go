package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) SaulGoodman(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Saul Goodman",
	})
}
