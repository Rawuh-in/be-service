package handler

import (
	"context"
	"encoding/json"
	"net/http"
	attendanceModel "rawuh-service/internal/attendance/model"
	attendanceService "rawuh-service/internal/attendance/service"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/middleware"
	"strconv"

	"github.com/gorilla/mux"
)

type AttendanceHandler struct {
	svc attendanceService.AttendanceService
}

func NewAttendanceHandler(svc attendanceService.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

// ListAttendance godoc
// @Summary List attendance records
// @Description Get paginated list of attendance records for an event
// @Tags attendance
// @Accept json
// @Produce json
// @Param page query int false "page"
// @Param limit query int false "limit"
// @Param sort query string false "sort field"
// @Param dir query string false "sort direction"
// @Param query query string false "search query"
// @Success 200 {object} attendanceModel.ListAttendanceResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/attendance/list [get]
func (h *AttendanceHandler) ListAttendance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &attendanceModel.ListAttendanceResponse{
		Error: false,
		Code:  http.StatusOK,
	}

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

	req := &attendanceModel.ListAttendanceRequest{
		Page:      int32(page),
		Limit:     int32(limit),
		Sort:      queryParams.Get("sort"),
		Dir:       queryParams.Get("dir"),
		Query:     queryParams.Get("query"),
		ProjectID: mux.Vars(r)["project_id"],
		EventId:   mux.Vars(r)["event_id"],
	}

	attendance, err := h.svc.ListAttendance(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	result.Data = attendance.Data

	utils.WriteJSONSuccess(w, result)
}

// GetAttendanceByID godoc
// @Summary Get attendance by ID
// @Description Get details for a specific attendance record
// @Tags attendance
// @Accept json
// @Produce json
// @Param attendance_id path string true "attendance id"
// @Success 200 {object} attendanceModel.GetAttendanceByIDResponse
// @Failure 404 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/attendance/{attendance_id} [get]
func (h *AttendanceHandler) GetAttendanceByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &attendanceModel.GetAttendanceByIDResponse{
		Error: false,
		Code:  http.StatusOK,
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	req := &attendanceModel.GetAttendanceByIDRequest{
		AttendanceID: mux.Vars(r)["attendance_id"],
		ProjectID:    mux.Vars(r)["project_id"],
		EventId:      mux.Vars(r)["event_id"],
	}

	attendance, err := h.svc.GetAttendanceByID(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	result.Data = attendance

	utils.WriteJSONSuccess(w, result)

}

// AddAttendanceBulk godoc
// @Summary Add attendance records in bulk
// @Description Create multiple attendance records for guests
// @Tags attendance
// @Accept json
// @Produce json
// @Param body body attendanceModel.CreateAttendanceRequest true "CreateAttendanceRequest"
// @Success 200 {object} attendanceModel.CreateAttendanceResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/attendance [post]
func (h *AttendanceHandler) AddAttendanceBulk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &attendanceModel.CreateAttendanceResponse{
		Error: false,
		Code:  http.StatusOK,
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var attendance attendanceModel.CreateAttendanceRequest
	if err := json.NewDecoder(r.Body).Decode(&attendance); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid Argument")
		return
	}

	req := &attendanceModel.CreateAttendanceRequest{
		ProjectID: mux.Vars(r)["project_id"],
		EventId:   mux.Vars(r)["event_id"],
		GuestID:   attendance.GuestID,
	}

	res, err := h.svc.AddAttendanceBulk(ctx, req)
	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	result.Data = res.Data

	utils.WriteJSONSuccess(w, result)
}

// CheckInOutAttendance godoc
// @Summary Check-in or check-out attendance
// @Description Bulk check-in or check-out attendance records
// @Tags attendance
// @Accept json
// @Produce json
// @Param body body attendanceModel.BulkCheckInOutAttendanceRequest true "BulkCheckInOutAttendanceRequest"
// @Success 200 {object} attendanceModel.BulkCheckInOutAttendanceResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/attendance/check [post]
func (h *AttendanceHandler) CheckInOutAttendance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &attendanceModel.BulkCheckInOutAttendanceResponse{
		Error: false,
		Code:  http.StatusOK,
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var attendanceReq attendanceModel.BulkCheckInOutAttendanceRequest
	if err := json.NewDecoder(r.Body).Decode(&attendanceReq); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid Argument")
		return
	}

	req := &attendanceModel.BulkCheckInOutAttendanceRequest{
		ProjectID:     mux.Vars(r)["project_id"],
		EventID:       mux.Vars(r)["event_id"],
		Type:          attendanceReq.Type,
		AttendanceIDs: attendanceReq.AttendanceIDs,
	}

	res, err := h.svc.CheckInOutAttendance(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	result.Data = res.Data

	utils.WriteJSONSuccess(w, result)
}

// DeleteAttendanceByID godoc
// @Summary Delete attendance records
// @Description Delete one or more attendance records by ID
// @Tags attendance
// @Accept json
// @Produce json
// @Param body body attendanceModel.DeleteAttendanceByIDRequest true "DeleteAttendanceByIDRequest"
// @Success 200 {object} attendanceModel.DeleteAttendanceByIDResponse
// @Failure 400 {object} utils.APIErrorResponse
// @Router /{project_id}/events/{event_id}/attendance [delete]
func (h *AttendanceHandler) DeleteAttendanceByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	result := &attendanceModel.DeleteAttendanceByIDResponse{
		Error: false,
		Code:  http.StatusOK,
	}

	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}

	var attendanceReq attendanceModel.DeleteAttendanceByIDRequest
	if err := json.NewDecoder(r.Body).Decode(&attendanceReq); err != nil {
		utils.WriteJSONError(w, http.StatusBadRequest, "Invalid Argument")
		return
	}

	req := &attendanceModel.DeleteAttendanceByIDRequest{
		ProjectID:     mux.Vars(r)["project_id"],
		EventID:       mux.Vars(r)["event_id"],
		AttendanceIDs: attendanceReq.AttendanceIDs,
	}

	res, err := h.svc.DeleteAttendanceByID(ctx, req)

	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}

	result.Data = res.Data

	utils.WriteJSONSuccess(w, result)
}
