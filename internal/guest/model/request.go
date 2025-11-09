package model

import "rawuh-service/internal/shared/model"

type ListGuestRequest struct {
	Page      int32  `json:"page"`
	Limit     int32  `json:"limit"`
	Sort      string `json:"sort"`
	Dir       string `json:"dir"`
	Query     string `json:"query"`
	EventId   string
	ProjectID string
}

type ListGuestResponse struct {
	Error      bool
	Code       int32
	Message    string
	Data       []*Guest
	Pagination *model.PaginationResponse
}

type CreateGuestRequest struct {
	ProjectID string
	Name      string
	Address   string
	Phone     string
	Email     string
	EventId   string
	EventData string
	GuestData string
}

type CreateGuestResponse struct {
	Error   bool
	Code    int32
	Message string
}
type UpdateGuestRequest struct {
	ProjectID string
	GuestID   string
	Name      string
	Address   string
	Phone     string
	Email     string
	EventId   string
	EventData string
	GuestData string
}

type UpdateGuestResponse struct {
	Error   bool
	Code    int32
	Message string
}

type GetGuestByIDRequest struct {
	ProjectID string
	GuestID   string
	EventId   string
}

type GetGuestByIDResponse struct {
	Error   bool
	Code    int32
	Message string
	Data    *Guest
}
type DeleteGuestByIDRequest struct {
	ProjectID string
	GuestID   string
	EventId   string
}

type DeleteGuestByIDResponse struct {
	Error   bool
	Code    int32
	Message string
}
