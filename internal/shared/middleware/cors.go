package middleware

import (
	"net/http"
	"strings"

	"rawuh-service/internal/shared/lib/utils"
)

// CORSMiddleware sets CORS headers for all routes and handles OPTIONS preflight.
// Uses ALLOWED_ORIGINS env (comma-separated) to control browser-accessible origins.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Check allowed origins from environment variable
		allowedOrigins := utils.GetEnv("ALLOWED_ORIGINS", "*")

		if allowedOrigins == "*" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else if origin != "" {
			// Check if origin is in allowed list
			allowed := strings.Split(allowedOrigins, ",")
			for _, allowedOrigin := range allowed {
				if strings.TrimSpace(allowedOrigin) == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", "true")
					break
				}
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Max-Age", "3600")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// For actual requests
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
