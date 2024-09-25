package middleware

import (
    "net/http"
)

func Authentication(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Authentication logic here
        next.ServeHTTP(w, r)
    })
}

func RateLimit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Rate limiting logic here
        next.ServeHTTP(w, r)
    })
}
