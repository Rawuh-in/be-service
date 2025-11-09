package middleware

import (
	"net/http"
	"strings"

	"rawuh-service/internal/shared/lib/utils"
)

// CORSMiddleware sets CORS headers for all routes and handles OPTIONS preflight.
// Uses ALLOWED_ORIGINS env (comma-separated) to control browser-accessible origins.
// Defaults to allowing http://localhost:3000 and http://127.0.0.1:3000 for local dev.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Build allowed origins list from env or use sensible defaults
		allowedEnv := utils.GetEnv("ALLOWED_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000")
		allowed := map[string]struct{}{}
		for _, o := range strings.Split(allowedEnv, ",") {
			if v := strings.TrimSpace(o); v != "" {
				allowed[v] = struct{}{}
			}
		}

		if origin == "" {
			// No Origin header (e.g., curl), keep permissive fallback
			w.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			// If origin is in allow-list, echo it; otherwise, omit header so browser blocks
			if _, ok := allowed[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// For preflight requests, respond with 200 OK directly if origin allowed
		if r.Method == http.MethodOptions {
			if origin == "" {
				w.WriteHeader(http.StatusOK)
				return
			}
			if _, ok := allowed[origin]; ok {
				w.WriteHeader(http.StatusOK)
				return
			}
			w.WriteHeader(http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
