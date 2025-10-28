package handler

import (
	"encoding/json"
	"net/http"
	projectModel "rawuh-service/internal/project/model"
	projectService "rawuh-service/internal/project/service"
	"rawuh-service/internal/shared/lib/utils"
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
func (h *ProjectHandler) ListProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &projectModel.ListProjectResponse{
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

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &projectModel.CreateProjectResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
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

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &projectModel.UpdateProjectResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
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

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &projectModel.DeleteProjectResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
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

func (h *ProjectHandler) DetailProject(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &projectModel.GetProjectDetailResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
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
