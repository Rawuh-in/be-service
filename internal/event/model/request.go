package model

import (
	"rawuh-service/internal/shared/model"
	"time"
)

type ListEventRequest struct {
	Page      int32  `json:"page"`
	Limit     int32  `json:"limit"`
	Sort      string `json:"sort"`
	Dir       string `json:"dir"`
	Query     string `json:"query"`
	ProjectID string
}

type ListEventResponse struct {
	Error      bool
	Code       int32
	Message    string
	Data       []*Event
	Pagination *model.PaginationResponse
}

type CreateEventRequest struct {
	EventName   string
	Description string
	Options     string
	StartDate   *time.Time
	EndDate     *time.Time
	UserID      string
	ProjectID   string
}

type CreateEventResponse struct {
	Error   bool
	Code    int32
	Message string
}

type UpdateEventRequest struct {
	ProjectID   string
	EventID     string
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
	EventsID  string
	ProjectID string
}

type DetailEventResponse struct {
	Error   bool
	Code    int32
	Message string
	Data    *Event
}

type DeleteEventRequest struct {
	EventsID  string
	ProjectID string
	UserID    string
}

type DeleteEventResponse struct {
	Error   bool
	Code    int32
	Message string
}
