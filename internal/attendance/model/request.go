package model

import "rawuh-service/internal/shared/model"

type ListAttendanceRequest struct {
	Page      int32  `json:"page"`
	Limit     int32  `json:"limit"`
	Sort      string `json:"sort"`
	Dir       string `json:"dir"`
	Query     string `json:"query"`
	EventId   string
	ProjectID string
}

type ListAttendanceResponse struct {
	Error      bool
	Code       int32
	Message    string
	Data       []*Attendance
	Pagination *model.PaginationResponse
}

type CreateAttendanceRequest struct {
	EventId   string
	ProjectID string
	GuestID   []*string
}

type GuestAttendanceResult struct {
	GuestID string `json:"guest_id"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type CreateAttendanceResponse struct {
	Error   bool                     `json:"error"`
	Code    int32                    `json:"code"`
	Message string                   `json:"message"`
	Data    []*GuestAttendanceResult `json:"results,omitempty"`
}

type UpdateAttendanceRequest struct {
	AttendanceID string
	EventId      string
	ProjectID    string
	GuestID      string
	Status       int64
	StatusStr    string
	Type         string
}

type UpdateAttendanceResponse struct {
	Error   bool
	Code    int32
	Message string
}

type GetAttendanceByIDRequest struct {
	ProjectID    string
	AttendanceID string
	EventId      string
}

type GetAttendanceByIDResponse struct {
	Error   bool
	Code    int32
	Message string
	Data    *Attendance
}

type DeleteAttendanceByIDRequest struct {
	ProjectID     string
	AttendanceIDs []string `json:"attendance_ids"`
	EventID       string
}

type AttendanceDeleteResult struct {
	AttendanceID string `json:"attendance_id"`
	Error        bool   `json:"error"`
	Message      string `json:"message"`
}

type DeleteAttendanceByIDResponse struct {
	Error   bool                      `json:"error"`
	Code    int32                     `json:"code"`
	Message string                    `json:"message"`
	Data    []*AttendanceDeleteResult `json:"data,omitempty"`
}

type BulkCheckInOutAttendanceRequest struct {
	ProjectID     string
	AttendanceIDs []string `json:"attendance_ids"`
	EventID       string
	Type          string
}

type BulkCheckInOutAttendanceResponse struct {
	Error   bool
	Code    int32
	Message string
	Data    []*CheckInOutResult
}

type CheckInOutResult struct {
	AttendanceID string `json:"attendance_id"`
	Error        bool   `json:"error"`
	Code         int32  `json:"code"`
	Message      string `json:"message"`
}
