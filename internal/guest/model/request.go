package guest_model

import "rawuh-service/internal/shared/model"

type ListProductRequest struct {
	Page  int32  `json:"page"`
	Limit int32  `json:"limit"`
	Sort  string `json:"sort"`
	Dir   string `json:"dir"`
	Query string `json:"query"`
}

type ListProductResponse struct {
	Error      bool
	Code       int32
	Message    string
	Data       []*Product
	Pagination *model.PaginationResponse
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Quantity    int32   `json:"quantity"`
}

type CreateProductResponse struct {
	Error   bool
	Code    int32
	Message string
}

type ListGuestRequest struct {
	Page    int32  `json:"page"`
	Limit   int32  `json:"limit"`
	Sort    string `json:"sort"`
	Dir     string `json:"dir"`
	Query   string `json:"query"`
	EventId string
}

type ListGuestResponse struct {
	Error      bool
	Code       int32
	Message    string
	Data       []*Guests
	Pagination *model.PaginationResponse
}

type CreateGuestRequest struct {
	Name    string
	Address string
	Phone   string
	Email   string
	EventId string
}

type CreateGuestResponse struct {
	Error   bool
	Code    int32
	Message string
}
