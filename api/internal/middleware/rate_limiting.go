package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	IPLimits map[string]*rate.Limiter // string is an IP
	mu       sync.Mutex
	interval time.Duration
	amount   int
}

func CreateRateLimiter(interval time.Duration, amount int) *IPRateLimiter {
	rateLimiter := &IPRateLimiter{
		mu:       sync.Mutex{},
		IPLimits: make(map[string]*rate.Limiter),
		interval: interval,
		amount:   amount,
	}

	return rateLimiter
}

func (rl *IPRateLimiter) Load(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	// Get the IP's rate limiter
	rlForIp, exists := rl.IPLimits[ip]
	if !exists {
		// Create a new rate limiter for this IP
		rlForIp = rate.NewLimiter(rate.Every(rl.interval), rl.amount)
		// Store it back into the map
		rl.IPLimits[ip] = rlForIp
	}

	// Return it
	return rl.IPLimits[ip]
}

func RateLimit(ipRateLimiter *IPRateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getIp(r)
		rateLimitForIp := ipRateLimiter.Load(ip)

		reservation := rateLimitForIp.Reserve()
		if reservation.Delay() > 0 {
			// Put the token back as reserve takes a token
			reservation.Cancel()
			retryAfter := math.Ceil(reservation.Delay().Seconds())
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter))
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]any{
				"message":             "Too many requests! Try again later",
				"retry_after_seconds": retryAfter,
			})
			return
		}

		next(w, r)
	}
}

func getIp(r *http.Request) string {
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

func (rl *IPRateLimiter) BackgroundCleanup(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 20)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.mu.Lock()
			// Get any buckets that are full and wipe them
			for k, v := range rl.IPLimits {
				if v.Tokens() >= float64(rl.amount) {
					delete(rl.IPLimits, k)
				}
			}
			rl.mu.Unlock()
		}
	}
}
