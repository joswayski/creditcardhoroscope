package middleware

import (
	"net/http"
)

func JSONHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodOptions {
			w.Header().Add("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}
