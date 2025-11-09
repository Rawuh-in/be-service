package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	eventModel "rawuh-service/internal/event/model"
	eventService "rawuh-service/internal/event/service"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/middleware"
	"strconv"

	"github.com/gorilla/mux"
)

type EventHandler struct {
	svc eventService.EventService
}

func NewEventHandler(svc eventService.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

// ListEvent godoc
// @Summary List events
// @Description Get paginated list of events for a project
// @Tags event
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} eventModel.ListEventResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/list [get]

func (h *EventHandler) ListEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.ListEventResponse{
		Error: false,
		Code:  http.StatusOK,
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

	req := &eventModel.ListEventRequest{
		Page:      int32(page),
		Limit:     int32(limit),
		Sort:      queryParams.Get("sort"),
		Dir:       queryParams.Get("dir"),
		Query:     queryParams.Get("query"),
		ProjectID: mux.Vars(r)["project_id"],
	}

	guests, err := h.svc.ListEvent(ctx, req)

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

// DetailEvent godoc
// @Summary Get event detail
// @Description Get details for a specific event
// @Tags event
// @Accept json
// @Produce json
// @Param event_id path string true "event id"
// @Success 200 {object} eventModel.DetailEventResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id} [get]

func (h *EventHandler) DetailEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.DetailEventResponse{
		Error: false,
		Code:  http.StatusOK,
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

	req := &eventModel.DetailEventRequest{
		EventsID:  mux.Vars(r)["event_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}

	guests, err := h.svc.DetailEvent(ctx, req)

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

// AddEvent godoc
// @Summary Create an event
// @Description Create a new event within a project
// @Tags event
// @Accept json
// @Produce json
// @Param body body eventModel.CreateEventRequest true "CreateEventRequest"
// @Success 200 {object} eventModel.CreateEventResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events [post]

func (h *EventHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.CreateEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success Create New Event",
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var p eventModel.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &eventModel.CreateEventRequest{
		EventName:    p.EventName,
		Description:  p.Description,
		EventOptions: p.EventOptions,
		GuestOptions: p.GuestOptions,
		StartDate:    p.StartDate,
		EndDate:      p.EndDate,
		UserID:       p.UserID,
		ProjectID:    mux.Vars(r)["project_id"],
	}

	if err := h.svc.AddEvent(ctx, req); err != nil {
		utils.HandleGrpcError(w, err)
		return

	}

	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update event details by ID
// @Tags event
// @Accept json
// @Produce json
// @Param event_id path string true "event id"
// @Param body body eventModel.UpdateEventRequest true "UpdateEventRequest"
// @Success 200 {object} eventModel.UpdateEventResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id} [put]

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.UpdateEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Success Update Event with id %s", mux.Vars(r)["event_id"]),
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var p eventModel.UpdateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &eventModel.UpdateEventRequest{
		ProjectID:    mux.Vars(r)["project_id"],
		EventID:      mux.Vars(r)["event_id"],
		EventName:    p.EventName,
		Description:  p.Description,
		EventOptions: p.EventOptions,
		GuestOptions: p.GuestOptions,
		StartDate:    p.StartDate,
		EndDate:      p.EndDate,
		UserID:       p.UserID,
	}

	if err := h.svc.UpdateEvent(ctx, req); err != nil {
		utils.HandleGrpcError(w, err)
		return

	}

	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Delete event by ID
// @Tags event
// @Accept json
// @Produce json
// @Param event_id path string true "event id"
// @Success 200 {object} eventModel.DeleteEventResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id} [delete]

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.DeleteEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Success Delete Event with id %s", mux.Vars(r)["event_id"]),
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	req := &eventModel.DeleteEventRequest{
		EventsID:  mux.Vars(r)["event_id"],
		UserID:    mux.Vars(r)["user_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}

	err := h.svc.DeleteEvent(ctx, req)

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
