package handler

import (
	"context"
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

func (h *AttendanceHandler) ListAttendance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// result := &attendanceModel.ListAttendanceResponse{
	// 	Error: false,
	// 	Code:  http.StatusOK,
	// }

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

	utils.WriteJSONSuccess(w, attendance)
}

func (h *AttendanceHandler) GetAttendanceByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// result := &attendanceModel.GetAttendanceByIDResponse{
	// 	Error: false,
	// 	Code:  http.StatusOK,
	// }

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

	utils.WriteJSONSuccess(w, attendance)

}

func (h *AttendanceHandler) CreateAttendance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// result := &attendanceModel.CreateAttendanceResponse{
	// 	Error: false,
	// 	Code:  http.StatusOK,
	// }
	if payloadMap, okp := middleware.GetAuthPayload(ctx); okp {
		ctx = context.WithValue(ctx, middleware.ContextKeyAuthPayload, payloadMap)
	}
	req := &attendanceModel.CreateAttendanceRequest{
		ProjectID: mux.Vars(r)["project_id"],
		EventId:   mux.Vars(r)["event_id"],
	}
	err := h.svc.CreateAttendance(ctx, req)
	if err != nil {
		utils.HandleGrpcError(w, err)
		return
	}
	utils.WriteJSONSuccess(w, nil)
}
