package model

import (
	"rawuh-service/internal/shared/model"
)

type ListProjectRequest struct {
	Page    int32  `json:"page"`
	Limit   int32  `json:"limit"`
	Sort    string `json:"sort"`
	Dir     string `json:"dir"`
	Query   string `json:"query"`
	EventId string
}

type ListProjectResponse struct {
	Error      bool
	Code       int32
	Message    string
	Data       []*Project
	Pagination *model.PaginationResponse
}

type CreateProjectRequest struct {
	ProjectName string
	UserID      string
}

type CreateProjectResponse struct {
	Error   bool
	Code    int32
	Message string
}

type UpdateProjectRequest struct {
	ProjectName string
	UserID      string
	ProjectID   string
	Status      int32
	StatusDesc  string
	Options     string
}

type UpdateProjectResponse struct {
	Error   bool
	Code    int32
	Message string
}

type DeleteProjectRequest struct {
	ProjectID string
}

type DeleteProjectResponse struct {
	Error   bool
	Code    int32
	Message string
}

type GetProjectDetailRequest struct {
	ProjectID string
}

type GetProjectDetailResponse struct {
	Error   bool
	Code    int32
	Message string
	Data    *Project
}
