package handler

import (
	"encoding/json"
	"net/http"
	eventModel "rawuh-service/internal/event/model"
	eventService "rawuh-service/internal/event/service"
	"rawuh-service/internal/shared/lib/utils"
	"strconv"

	"github.com/gorilla/mux"
)

type EventHandler struct {
	svc eventService.EventService
}

func NewEventHandler(svc eventService.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

func (h *EventHandler) ListEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.ListEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
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

func (h *EventHandler) DetailEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.DetailEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
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

func (h *EventHandler) AddEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.CreateEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success Create New Event",
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
		EventName:   p.EventName,
		Description: p.Description,
		Options:     p.Options,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		UserID:      p.UserID,
		ProjectID:   mux.Vars(r)["project_id"],
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

func (h *EventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.UpdateEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success Create New Event",
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
		ProjectID:   mux.Vars(r)["project_id"],
		EventID:     mux.Vars(r)["event_id"],
		EventName:   p.EventName,
		Description: p.Description,
		Options:     p.Options,
		StartDate:   p.StartDate,
		EndDate:     p.EndDate,
		UserID:      p.UserID,
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

func (h *EventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &eventModel.DeleteEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}

	req := &eventModel.DeleteEventRequest{
		EventsID:  mux.Vars(r)["event_id"],
		UserID:    mux.Vars(r)["user_id"],
		ProjectID: mux.Vars(r)["project_id"],
	}

	event, err := h.svc.DeleteEvent(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}
	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}
