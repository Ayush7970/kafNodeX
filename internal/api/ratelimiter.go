// internal/api/ratelimiter.go
package api

import (
    "net/http"
    "golang.org/x/time/rate"
)

// Rate limiter settings
const (
    globalRequestsPerSecond = 5  // Set a global rate limit
    globalBurstLimit        = 10  // Allow up to 2 requests in a burst
)

// Global rate limiter instance
var globalLimiter = rate.NewLimiter(rate.Limit(globalRequestsPerSecond), globalBurstLimit)

// RateLimiterMiddleware applies a global rate limit for all requests
func RateLimiterMiddleware() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            log.Println("RateLimiterMiddleware invoked")  // Log every request

            // Check if the request should be allowed by the global limiter
            if !globalLimiter.Allow() {
                log.Println("Rate limit exceeded")  // Log when the rate limit is triggered
                http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
                return
            }

            log.Println("Request allowed")  // Log allowed requests
            next.ServeHTTP(w, r)
        })
    }
}
