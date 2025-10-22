package guest_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	guestModel "rawuh-service/internal/guest/model"
	guestService "rawuh-service/internal/guest/service"
	"rawuh-service/internal/shared/lib/utils"

	"github.com/gorilla/mux"
)

type GuestHandler struct {
	svc guestService.GuestService
}

func NewGuestHandler(svc guestService.GuestService) *GuestHandler {
	return &GuestHandler{svc: svc}
}

func (h *GuestHandler) AddGuest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &guestModel.CreateGuestResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success Create New Guest",
	}

	var p guestModel.CreateGuestRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &guestModel.CreateGuestRequest{
		Name:    p.Name,
		Address: p.Address,
		Phone:   p.Phone,
		Email:   p.Email,
		EventId: p.EventId,
	}
	if err := h.svc.AddGuest(ctx, req); err != nil {
		utils.HandleGrpcError(w, err)
		return

	}

	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *GuestHandler) UpdateGuestByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &guestModel.UpdateGuestResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success Update Guest",
	}

	var p guestModel.Guests
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &guestModel.UpdateGuestRequest{
		EventId: mux.Vars(r)["event_id"],
		GuestID: mux.Vars(r)["guest_id"],
		Name:    p.Name,
		Address: p.Address,
		Phone:   p.Phone,
		Email:   p.Email,
	}
	if err := h.svc.UpdateGuestByID(ctx, req); err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func (h *GuestHandler) ListGuests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &guestModel.ListGuestResponse{
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

	req := &guestModel.ListGuestRequest{
		Page:    int32(page),
		Limit:   int32(limit),
		Sort:    queryParams.Get("sort"),
		Dir:     queryParams.Get("dir"),
		Query:   queryParams.Get("query"),
		EventId: mux.Vars(r)["event_id"],
	}

	guests, err := h.svc.ListGuests(ctx, req)

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

func (h *GuestHandler) GetGuestByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &guestModel.GetGuestByIDResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}

	req := &guestModel.GetGuestByIDRequest{
		EventId: mux.Vars(r)["event_id"],
		GuestID: mux.Vars(r)["guest_id"],
	}

	guest, err := h.svc.GetGuestByID(ctx, req)
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

func (h *GuestHandler) DeleteGuestByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &guestModel.GetGuestByIDResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}

	req := &guestModel.DeleteGuestByIDRequest{
		EventId: mux.Vars(r)["event_id"],
		GuestID: mux.Vars(r)["guest_id"],
	}

	guest, err := h.svc.DeleteGuestByID(ctx, req)
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
