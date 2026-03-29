package server

import "net/http"

func (s *Server) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
}
