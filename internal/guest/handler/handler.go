package guest_handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	guest_model "rawuh-service/internal/guest/model"
	guest_service "rawuh-service/internal/guest/service"
	"rawuh-service/internal/shared/lib/utils"
)

type GuestHandler struct {
	svc guest_service.GuestService
}

func NewGuestHandler(svc guest_service.GuestService) *GuestHandler {
	return &GuestHandler{svc: svc}
}

// func (h *GuestHandler) AddProduct(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	result := &guest_model.CreateProductResponse{
// 		Error:   false,
// 		Code:    http.StatusOK,
// 		Message: "Success Create Product",
// 	}

// 	var p guest_model.Product
// 	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
// 		result.Error = true
// 		result.Code = http.StatusInternalServerError
// 		result.Message = "Invalid Argument"
// 		w.Header().Add("content-type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(result)
// 		return
// 	}

// 	req := &guest_model.CreateProductRequest{
// 		Name:        p.Name,
// 		Price:       p.Price,
// 		Description: p.Description,
// 		Quantity:    p.Quantity,
// 	}
// 	if err := h.svc.AddProduct(ctx, req); err != nil {
// 		result.Error = true
// 		result.Code = http.StatusForbidden
// 		result.Message = err.Error()
// 		w.Header().Add("content-type", "application/json")
// 		w.WriteHeader(http.StatusBadRequest)
// 		json.NewEncoder(w).Encode(result)
// 		return
// 	}

// 	result.Error = false
// 	result.Code = http.StatusOK
// 	w.Header().Add("content-type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(result)
// }

// func (h *GuestHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
// 	ctx := r.Context()

// 	result := &guest_model.ListProductResponse{
// 		Error:   false,
// 		Code:    http.StatusOK,
// 		Message: "Success",
// 	}

// 	queryParams := r.URL.Query()

// 	page, _ := strconv.Atoi(queryParams.Get("page"))
// 	limit, _ := strconv.Atoi(queryParams.Get("limit"))

// 	if page <= 0 {
// 		page = 1
// 	}
// 	if limit <= 0 {
// 		limit = 10
// 	}

// 	req := &guest_model.ListProductRequest{
// 		Page:  int32(page),
// 		Limit: int32(limit),
// 		Sort:  queryParams.Get("sort"),
// 		Dir:   queryParams.Get("dir"),
// 		Query: queryParams.Get("query"),
// 	}

// 	products, err := h.svc.ListProducts(ctx, req)

// 	if err != nil {
// 		result.Error = true
// 		result.Code = http.StatusInternalServerError
// 		result.Message = err.Error()
// 		w.Header().Add("content-type", "application/json")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		json.NewEncoder(w).Encode(result)
// 		return
// 	}

//		result.Error = false
//		result.Code = http.StatusOK
//		w.Header().Add("content-type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode(products)
//	}
func (h *GuestHandler) AddGuest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &guest_model.CreateProductResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success Create Product",
	}

	var p guest_model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &guest_model.CreateGuestRequest{
		Name: p.Name,
	}
	if err := h.svc.AddGuest(ctx, req); err != nil {
		result.Error = true
		result.Code = http.StatusForbidden
		result.Message = err.Error()
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
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

	result := &guest_model.ListProductResponse{
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

	req := &guest_model.ListGuestRequest{
		Page:    int32(page),
		Limit:   int32(limit),
		Sort:    queryParams.Get("sort"),
		Dir:     queryParams.Get("dir"),
		Query:   queryParams.Get("query"),
		EventId: queryParams.Get("event_id"),
	}

	products, err := h.svc.ListGuests(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}
	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}
