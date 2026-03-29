package middleware

import (
	"net/http"
	"slices"
)

var allowedOrigins = []string{"http://localhost:5173", "http://localhost:3000", "https://creditcardhoroscope.com", "https://staging.creditcardhoroscope.com"}

func CORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originDomain := r.Header.Get("Origin")
		if slices.Contains(allowedOrigins, originDomain) {
			w.Header().Set("Access-Control-Allow-Origin", originDomain)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue
		next.ServeHTTP(w, r)
	})
}
