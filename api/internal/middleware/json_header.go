package middleware

import (
	"net/http"
	"slices"
)

var allowedMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

func JSONHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if slices.Contains(allowedMethods, r.Method) {
			w.Header().Add("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}
