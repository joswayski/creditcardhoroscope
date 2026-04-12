package middleware

import (
	"net/http"
)

func BodySize(next http.HandlerFunc, maxSize uint64) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, int64(maxSize))
		next.ServeHTTP(w, r)
	}
}
