package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	redisPkg "rawuh-service/internal/shared/redis"
)

type ContextKey string

const ContextKeyAuthPayload ContextKey = "auth_payload"

func AuthMiddleware(rdb *redisPkg.Redis) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth != "" {
				// Remove "Bearer " prefix if it exists
				token := auth
				if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
					token = strings.TrimSpace(auth[len("Bearer "):])
				} else {
					token = strings.TrimSpace(auth)
				}
				if token != "" {
					key := "access_token:" + token
					if val, err := rdb.Get(r.Context(), key); err == nil && val != "" {
						var payload map[string]interface{}
						if err := json.Unmarshal([]byte(val), &payload); err == nil {
							ctx := context.WithValue(r.Context(), ContextKeyAuthPayload, payload)
							r = r.WithContext(ctx)
						}
					}
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if v := r.Context().Value(ContextKeyAuthPayload); v == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"error": true, "message": "unauthenticated"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetAuthPayload(ctx context.Context) (map[string]interface{}, bool) {
	v := ctx.Value(ContextKeyAuthPayload)
	if v == nil {
		return nil, false
	}
	payload, ok := v.(map[string]interface{})
	return payload, ok
}

func GetStringClaim(payload map[string]interface{}, key string) (string, bool) {
	if payload == nil {
		return "", false
	}
	if v, ok := payload[key]; ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}

func GetInt64Claim(payload map[string]interface{}, key string) (int64, bool) {
	if payload == nil {
		return 0, false
	}
	if v, ok := payload[key]; ok {
		switch t := v.(type) {
		case float64:
			return int64(t), true
		case int:
			return int64(t), true
		case int64:
			return t, true
		case string:
			if parsed, err := strconv.ParseInt(t, 10, 64); err == nil {
				return parsed, true
			}
		}
	}
	return 0, false
}

type AuthClaims struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	UserID    int64  `json:"user_id"`
	ProjectID int64  `json:"project_id"`
	EventID   int64  `json:"event_id"`
	UserType  string `json:"usertype"`
}

func ParseAuthClaims(payload map[string]interface{}) (AuthClaims, bool) {
	if payload == nil {
		return AuthClaims{}, false
	}
	var c AuthClaims
	if s, ok := GetStringClaim(payload, "username"); ok {
		c.Username = s
	}
	if s, ok := GetStringClaim(payload, "name"); ok {
		c.Name = s
	}
	if v, ok := GetInt64Claim(payload, "user_id"); ok {
		c.UserID = v
	}
	if v, ok := GetInt64Claim(payload, "project_id"); ok {
		c.ProjectID = v
	}
	if v, ok := GetInt64Claim(payload, "event_id"); ok {
		c.EventID = v
	}
	if s, ok := GetStringClaim(payload, "usertype"); ok {
		c.UserType = s
	}
	return c, true
}

func GetAuthClaimsFromContext(ctx context.Context) (AuthClaims, bool) {
	if payload, ok := GetAuthPayload(ctx); ok {
		return ParseAuthClaims(payload)
	}
	return AuthClaims{}, false
}
