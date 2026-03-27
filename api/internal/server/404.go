package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) FourOhFour(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Not Found",
	})
}
