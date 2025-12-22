package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	guestModel "rawuh-service/internal/guest/model"
	guestService "rawuh-service/internal/guest/service"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/middleware"

	"github.com/gorilla/mux"
)

type GuestHandler struct {
	svc guestService.GuestService
}

func NewGuestHandler(svc guestService.GuestService) *GuestHandler {
	return &GuestHandler{svc: svc}
}

// AddGuest godoc
// @Summary Create a new guest
// @Description Create guest for an event
// @Tags guest
// @Accept json
// @Produce json
// @Param body body guestModel.CreateGuestRequest true "CreateGuestRequest"
// @Success 200 {object} guestModel.CreateGuestResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/guests [post]

func (h *GuestHandler) AddGuest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &guestModel.CreateGuestResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success create new guest",
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var guestReq guestModel.CreateGuestRequest
	if err := json.NewDecoder(r.Body).Decode(&guestReq); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid Argument")
		return
	}

	req := &guestModel.CreateGuestRequest{
		Name:      guestReq.Name,
		Address:   guestReq.Address,
		Phone:     guestReq.Phone,
		Email:     guestReq.Email,
		EventData: guestReq.EventData,
		GuestData: guestReq.GuestData,
		EventId:   mux.Vars(r)["event_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}
	if err := h.svc.AddGuest(ctx, req); err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	utils.WriteJSONSuccess(w, result)
}

// UpdateGuestByID godoc
// @Summary Update guest by ID
// @Description Update guest details
// @Tags guest
// @Accept json
// @Produce json
// @Param body body guestModel.UpdateGuestRequest true "UpdateGuestRequest"
// @Success 200 {object} guestModel.UpdateGuestResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/guests/{guest_id} [put]

func (h *GuestHandler) UpdateGuestByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &guestModel.UpdateGuestResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Success Update Guest with id %s ", mux.Vars(r)["guest_id"]),
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var guestReq guestModel.UpdateGuestRequest
	if err := json.NewDecoder(r.Body).Decode(&guestReq); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid Argument")
		return
	}

	req := &guestModel.UpdateGuestRequest{
		ProjectID: mux.Vars(r)["project_id"],
		EventId:   mux.Vars(r)["event_id"],
		GuestID:   mux.Vars(r)["guest_id"],
		EventData: guestReq.EventData,
		GuestData: guestReq.GuestData,
		Name:      guestReq.Name,
		Address:   guestReq.Address,
		Phone:     guestReq.Phone,
		Email:     guestReq.Email,
	}
	if err := h.svc.UpdateGuestByID(ctx, req); err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	utils.WriteJSONSuccess(w, result)
}

// ListGuests godoc
// @Summary List guests
// @Description Get list of guests for an event
// @Tags guest
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} guestModel.ListGuestResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/guests/list [get]

func (h *GuestHandler) ListGuests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// result := &guestModel.ListGuestResponse{
	// 	Error: false,
	// 	Code:  http.StatusOK,
	// }

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

	req := &guestModel.ListGuestRequest{
		Page:      int32(page),
		Limit:     int32(limit),
		Sort:      queryParams.Get("sort"),
		Dir:       queryParams.Get("dir"),
		Query:     queryParams.Get("query"),
		EventId:   mux.Vars(r)["event_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}

	guests, err := h.svc.ListGuests(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	utils.WriteJSONSuccess(w, guests)
}

// GetGuestByID godoc
// @Summary Get guest by ID
// @Description Get guest details
// @Tags guest
// @Accept json
// @Produce json
// @Param guest_id path string true "guest id"
// @Success 200 {object} guestModel.GetGuestByIDResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/guests/{guest_id} [get]

func (h *GuestHandler) GetGuestByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// result := &guestModel.GetGuestByIDResponse{
	// 	Error: false,
	// 	Code:  http.StatusOK,
	// }

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	req := &guestModel.GetGuestByIDRequest{
		EventId:   mux.Vars(r)["event_id"],
		GuestID:   mux.Vars(r)["guest_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}

	guest, err := h.svc.GetGuestByID(ctx, req)
	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	utils.WriteJSONSuccess(w, guest)
}

// DeleteGuestByID godoc
// @Summary Delete guest by ID
// @Description Delete a guest
// @Tags guest
// @Accept json
// @Produce json
// @Param guest_id path string true "guest id"
// @Success 200 {object} guestModel.GetGuestByIDResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/guests/{guest_id} [delete]

func (h *GuestHandler) DeleteGuestByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &guestModel.GetGuestByIDResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	req := &guestModel.DeleteGuestByIDRequest{
		EventId:   mux.Vars(r)["event_id"],
		GuestID:   mux.Vars(r)["guest_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}

	err := h.svc.DeleteGuestByID(ctx, req)
	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	utils.WriteJSONSuccess(w, result)
}
