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

		// Set standard CORS headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle Origin
		if origin != "" {
			// Always allow localhost origins in development
			if strings.Contains(origin, "localhost") || strings.Contains(origin, "127.0.0.1") {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				// For production, check against allowed origins
				allowedEnv := utils.GetEnv("ALLOWED_ORIGINS", "*")
				if allowedEnv == "*" {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else {
					allowed := strings.Split(allowedEnv, ",")
					for _, o := range allowed {
						if strings.TrimSpace(o) == origin {
							w.Header().Set("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			}
		}

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			if origin != "" {
				w.Header().Set("Content-Type", "text/plain")
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}

		// For actual requests, set content type to JSON
		if r.Method != http.MethodOptions {
			w.Header().Set("Content-Type", "application/json")
		}

		next.ServeHTTP(w, r)
	})
}
