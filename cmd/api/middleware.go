package main

import (
    "net/http"
    "os"
)

// Middleware function to check for the API key in each request
func apiKeyAuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        apiKey := r.Header.Get("X-API-Key")
        expectedApiKey := os.Getenv("API_KEY") // Retrieve the expected API key from environment variables

        // Check if the API key is missing or incorrect
        if apiKey == "" || apiKey != expectedApiKey {
            http.Error(w, "Unauthorized: Invalid API key", http.StatusUnauthorized)
            return
        }

        // If the API key is valid, continue with the request
        next.ServeHTTP(w, r)
    })
}
