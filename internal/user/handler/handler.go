package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/middleware"
	userModel "rawuh-service/internal/user/model"
	userService "rawuh-service/internal/user/service"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	svc userService.UserService
}

func NewUserHandler(svc userService.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// AddUser godoc
// @Summary Create a new user
// @Description Create user and create auth row when password provided
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param body body userModel.CreateUserRequest true "CreateUserRequest"
// @Success 200 {object} userModel.CreateUserResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /users [post]

func (h *UserHandler) AddUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &userModel.CreateUserResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success Create New Guest",
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var p userModel.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &userModel.CreateUserRequest{
		ProjectID: p.ProjectID,
		EventId:   p.EventId,
		Name:      p.Name,
		Username:  p.Username,
		Password:  p.Password,
		UserType:  p.UserType,
		Email:     p.Email,
	}
	if err := h.svc.AddUser(ctx, req); err != nil {
		utils.HandleGrpcError(w, err)
		return

	}

	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// UpdateUserByID godoc
// @Summary Update a user
// @Description Update user details by ID
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param user_id path string true "user id"
// @Param body body userModel.UpdateUserRequest true "UpdateUserRequest"
// @Success 200 {object} userModel.UpdateUserResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /users/{user_id} [put]

func (h *UserHandler) UpdateUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &userModel.UpdateUserResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success Update Guest",
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var p userModel.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &userModel.UpdateUserRequest{
		ProjectID: mux.Vars(r)["project_id"],
		EventId:   mux.Vars(r)["event_id"],
		UserID:    mux.Vars(r)["user_id"],
		Name:      p.Name,
		Username:  p.Username,
		UserType:  p.UserType,
		Email:     p.Email,
	}
	if err := h.svc.UpdateUserByID(ctx, req); err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// ListUsers godoc
// @Summary List users
// @Description Get paginated list of users
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} userModel.ListUserResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /users/list [get]

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &userModel.ListUserResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	queryParams := r.URL.Query()

	page, _ := strconv.Atoi(queryParams.Get("page"))
	limit, _ := strconv.Atoi(queryParams.Get("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	req := &userModel.ListUserRequest{
		Page:      int32(page),
		Limit:     int32(limit),
		Sort:      queryParams.Get("sort"),
		Dir:       queryParams.Get("dir"),
		Query:     queryParams.Get("query"),
		EventId:   mux.Vars(r)["event_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}

	guests, err := h.svc.ListUsers(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}
	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guests)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user details by ID
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param user_id path string true "user id"
// @Success 200 {object} userModel.GetUserByIDResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /users/{user_id} [get]

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &userModel.GetUserByIDResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	req := &userModel.GetUserByIDRequest{
		EventId:   mux.Vars(r)["event_id"],
		UserID:    mux.Vars(r)["user_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}

	guest, err := h.svc.GetUserByID(ctx, req)
	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}
	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guest)
}

// DeleteUserByID godoc
// @Summary Delete user by ID
// @Description Delete a user
// @Tags user
// @Accept json
// @Produce json
// @Security Bearer
// @Param user_id path string true "user id"
// @Success 200 {object} userModel.GetUserByIDResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /users/{user_id} [delete]

func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &userModel.GetUserByIDResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	req := &userModel.DeleteUserByIDRequest{
		EventId:   mux.Vars(r)["event_id"],
		UserID:    mux.Vars(r)["user_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}

	err := h.svc.DeleteUserByID(ctx, req)
	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
