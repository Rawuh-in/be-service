package handler

import (
	"context"
	"encoding/json"
	"net/http"
	projectModel "rawuh-service/internal/project/model"
	projectService "rawuh-service/internal/project/service"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/middleware"
	"strconv"

	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	svc projectService.ProjectService
}

func NewProjectHandler(svc projectService.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		svc: svc,
	}
}

// ListProject godoc
// @Summary List projects
// @Description Get paginated list of projects
// @Tags project
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Success 200 {object} projectModel.ListProjectResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /project/list [get]

func (h *ProjectHandler) ListProject(w http.ResponseWriter, r *http.Request) {
	result := &projectModel.ListProjectResponse{
		Error: false,
		Code:  http.StatusOK,
	}

	ctx := r.Context()

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
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

	req := &projectModel.ListProjectRequest{
		Page:    int32(page),
		Limit:   int32(limit),
		Sort:    queryParams.Get("sort"),
		Dir:     queryParams.Get("dir"),
		Query:   queryParams.Get("query"),
		EventId: mux.Vars(r)["event_id"],
	}

	guests, err := h.svc.ListProjects(ctx, req)

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

// CreateProject godoc
// @Summary Create a project
// @Description Create a new project
// @Tags project
// @Accept json
// @Produce json
// @Param body body projectModel.CreateProjectRequest true "CreateProjectRequest"
// @Success 200 {object} projectModel.CreateProjectResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /project [post]

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &projectModel.CreateProjectResponse{
		Error: false,
		Code:  http.StatusOK,
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var p projectModel.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &projectModel.CreateProjectRequest{
		ProjectName: p.ProjectName,
		UserID:      p.UserID,
	}

	err := h.svc.CreateProject(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}
	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update project by ID
// @Tags project
// @Accept json
// @Produce json
// @Param project_id path string true "project id"
// @Param body body projectModel.UpdateProjectRequest true "UpdateProjectRequest"
// @Success 200 {object} projectModel.UpdateProjectResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /project/{project_id} [put]

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &projectModel.UpdateProjectResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var p projectModel.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &projectModel.UpdateProjectRequest{
		ProjectID:   mux.Vars(r)["project_id"],
		ProjectName: p.ProjectName,
		UserID:      p.UserID,
	}

	err := h.svc.UpdateProject(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}
	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete project by ID
// @Tags project
// @Accept json
// @Produce json
// @Param project_id path string true "project id"
// @Success 200 {object} projectModel.DeleteProjectResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /project/{project_id} [delete]

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &projectModel.DeleteProjectResponse{
		Error: false,
		Code:  http.StatusOK,
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var p projectModel.DeleteProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		result.Error = true
		result.Code = http.StatusInternalServerError
		result.Message = "Invalid Argument"
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(result)
		return
	}

	req := &projectModel.DeleteProjectRequest{
		ProjectID: mux.Vars(r)["project_id"],
	}

	err := h.svc.DeleteProject(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}
	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

// DetailProject godoc
// @Summary Get project detail
// @Description Retrieve project details by ID
// @Tags project
// @Accept json
// @Produce json
// @Param project_id path string true "project id"
// @Success 200 {object} projectModel.GetProjectDetailResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /project/{project_id} [get]

func (h *ProjectHandler) DetailProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &projectModel.GetProjectDetailResponse{
		Error: false,
		Code:  http.StatusOK,
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	req := &projectModel.GetProjectDetailRequest{
		ProjectID: mux.Vars(r)["project_id"],
	}

	project, err := h.svc.GetProjectDetail(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}
	result.Error = false
	result.Code = http.StatusOK
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(project)
}
