package event_model

import (
	"rawuh-service/internal/shared/model"
	"time"
)

type ListEventRequest struct {
	Page   int32  `json:"page"`
	Limit  int32  `json:"limit"`
	Sort   string `json:"sort"`
	Dir    string `json:"dir"`
	Query  string `json:"query"`
	UserID string
}

type ListEventResponse struct {
	Error      bool
	Code       int32
	Message    string
	Data       []*Events
	Pagination *model.PaginationResponse
}

type CreateEventRequest struct {
	EventName   string
	Description string
	Options     string
	StartDate   *time.Time
	EndDate     *time.Time
	UserID      string
}

type CreateEventResponse struct {
	Error   bool
	Code    int32
	Message string
}

type UpdateEventRequest struct {
	EventsID    string
	EventName   string
	Description string
	Options     string
	StartDate   *time.Time
	EndDate     *time.Time
	UserID      string
}

type UpdateEventResponse struct {
	Error   bool
	Code    int32
	Message string
}

type DetailEventRequest struct {
	EventsID string
	UserID   string
}

type DetailEventResponse struct {
	Error   bool
	Code    int32
	Message string
	Data    *Events
}

type DeleteEventRequest struct {
	EventsID string
	UserID   string
}

type DeleteEventResponse struct {
	Error   bool
	Code    int32
	Message string
}
