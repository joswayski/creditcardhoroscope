package middleware

import (
	"encoding/json"
	"net/http"
)

func IPWhitelist(next http.HandlerFunc, ipWhitelist map[string]bool, environment string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := GetIp(r)
		isLocal := ip == "::1" && environment == "development"
		if !ipWhitelist[ip] && !isLocal { // local stripe
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{})
			return
		}
		next.ServeHTTP(w, r)
	}

}
