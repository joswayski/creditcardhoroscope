package middleware

import (
	"net/http"
)

func BodySize(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		next.ServeHTTP(w, r)
	})
}
