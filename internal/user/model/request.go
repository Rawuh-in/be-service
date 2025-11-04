package model

import "rawuh-service/internal/shared/model"

type ListUserRequest struct {
	Page      int32  `json:"page"`
	Limit     int32  `json:"limit"`
	Sort      string `json:"sort"`
	Dir       string `json:"dir"`
	Query     string `json:"query"`
	EventId   string
	ProjectID string
}

type ListUserResponse struct {
	Error      bool
	Code       int32
	Message    string
	Data       []*User
	Pagination *model.PaginationResponse
}

type CreateUserRequest struct {
	ProjectID string
	Name      string
	Username  string
	Password  string
	UserType  string
	Email     string
	UserID    string
	EventId   string
}

type CreateUserResponse struct {
	Error   bool
	Code    int32
	Message string
}
type UpdateUserRequest struct {
	ProjectID string
	Name      string
	Username  string
	UserType  string
	Email     string
	UserID    string
	EventId   string
}

type UpdateUserResponse struct {
	Error   bool
	Code    int32
	Message string
}

type GetUserByIDRequest struct {
	ProjectID string
	UserID    string
	EventId   string
}

type GetUserByIDResponse struct {
	Error   bool
	Code    int32
	Message string
	Data    *User
}
type DeleteUserByIDRequest struct {
	ProjectID string
	UserID    string
	EventId   string
}

type DeleteUserByIDResponse struct {
	Error   bool
	Code    int32
	Message string
}
