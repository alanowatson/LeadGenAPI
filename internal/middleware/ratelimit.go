package middleware

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	"github.com/alanowatson/LeadGenAPI/pkg/util"
)

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(2, 5) // 2 requests per second, burst of 5
		visitors[ip] = &visitor{limiter, time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

func RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		limiter := getVisitor(ip)
		if !limiter.Allow() {
			util.RespondWithError(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		next.ServeHTTP(w, r)
	}
}

func CleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		mu.Lock()
		for ip, v := range visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(visitors, ip)
			}
		}
		mu.Unlock()
	}
}
