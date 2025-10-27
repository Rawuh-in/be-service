package model

import "rawuh-service/internal/shared/model"

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
	Data       []*Guest
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
type UpdateGuestRequest struct {
	GuestID string
	Name    string
	Address string
	Phone   string
	Email   string
	EventId string
}

type UpdateGuestResponse struct {
	Error   bool
	Code    int32
	Message string
}

type GetGuestByIDRequest struct {
	GuestID string
	EventId string
}

type GetGuestByIDResponse struct {
	Error   bool
	Code    int32
	Message string
	Data    *Guest
}
type DeleteGuestByIDRequest struct {
	GuestID string
	EventId string
}

type DeleteGuestByIDResponse struct {
	Error   bool
	Code    int32
	Message string
}
