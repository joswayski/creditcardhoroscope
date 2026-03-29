package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hi! You probably meant to go to one of the other routes. Make sure to check the documentation!",
	})

}
