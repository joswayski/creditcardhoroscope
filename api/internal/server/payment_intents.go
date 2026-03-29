package server

import "net/http"

func (s *Server) CreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
