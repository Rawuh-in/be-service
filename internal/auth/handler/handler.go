package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	authService "rawuh-service/internal/auth/service"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/logger"
	"rawuh-service/internal/shared/redis"
	userDb "rawuh-service/internal/user/repository"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authSvc authService.AuthService
	userDb  *userDb.UserRepository
	rdb     *redis.Redis
	logger  *logger.Logger
}

func NewAuthHandler(a authService.AuthService, u *userDb.UserRepository, r *redis.Redis, l *logger.Logger) *AuthHandler {
	return &AuthHandler{authSvc: a, userDb: u, rdb: r, logger: l}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Error       bool   `json:"error"`
	Code        int    `json:"code"`
	AccessToken string `json:"access_token"`
	Message     string `json:"message"`
}

// Login godoc
// @Summary Login with username and password
// @Description Authenticate user and return an access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body loginRequest true "Login credentials"
// @Success 200 {object} loginResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Failure 401 {object} utils.APIErrorResponse
// @Router /login [post]

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Support Basic auth header first: Authorization: Basic base64(username:password)
	username, password, ok := r.BasicAuth()
	var req loginRequest
	if ok {
		req.Username = username
		req.Password = password
	} else {
		// fallback to JSON body
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.HandleGrpcError(w, err)
			return
		}
	}

	authRow, err := h.authSvc.Authenticate(ctx, req.Username, req.Password)
	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	// get user info
	userIDStr := strconv.FormatInt(authRow.UserID, 10)
	user, err := h.userDb.GetUserByID(ctx, userIDStr)
	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	// prepare token payload
	payload := map[string]interface{}{
		"username":   authRow.Username,
		"name":       user.Name,
		"user_id":    authRow.UserID,
		"project_id": authRow.ProjectID,
		"event_id":   user.EventId,
		"usertype":   user.UserType,
	}

	// generate token
	token := uuid.New().String()
	key := "access_token:" + token
	// store in redis, 24h
	if err := h.rdb.Set(ctx, key, payload, 24*time.Hour); err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	res := &loginResponse{
		Error:       false,
		Code:        http.StatusOK,
		AccessToken: token,
		Message:     "success",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// TokenInfo returns the payload stored for the provided Bearer token.
// It reads the Authorization: Bearer <token> header, fetches the payload
// from Redis and returns it as JSON.
func (h *AuthHandler) TokenInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	auth := r.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": true, "message": "missing bearer token"})
		return
	}

	token := strings.TrimSpace(auth[len("Bearer "):])
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": true, "message": "invalid token"})
		return
	}

	key := "access_token:" + token
	val, err := h.rdb.Get(ctx, key)
	if err != nil || val == "" {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{"error": true, "message": "token not found"})
		return
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(val), &payload); err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}
