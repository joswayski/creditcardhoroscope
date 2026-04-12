package middleware

import (
	"encoding/json"
	"net/http"
)

func IPWhitelist(next http.HandlerFunc, ipWhitelist map[string]bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := GetIp(r)
		if !ipWhitelist[ip] && ip != "::1" { // local stripe
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}
		next.ServeHTTP(w, r)
	}

}
