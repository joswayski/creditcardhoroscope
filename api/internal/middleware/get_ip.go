package middleware

import (
	"net"
	"net/http"
	"strings"
)

func GetIp(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		ip = strings.Split(ip, ",")[0] // can be multiple
	}

	if ip == "" {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			ip = host
		} else {
			ip = r.RemoteAddr
		}
	}

	return ip
}
