package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	attendanceModel "rawuh-service/internal/attendance/model"
	attendanceDb "rawuh-service/internal/attendance/repository"
	guestModel "rawuh-service/internal/guest/model"
	"rawuh-service/internal/shared/constant"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/logger"
	"rawuh-service/internal/shared/middleware"
	"rawuh-service/internal/shared/model"
	"strconv"
	"strings"

	"go.elastic.co/apm/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AttendanceService interface {
	ListAttendance(ctx context.Context, req *attendanceModel.ListAttendanceRequest) (*attendanceModel.ListAttendanceResponse, error)
	GetAttendanceByID(ctx context.Context, req *attendanceModel.GetAttendanceByIDRequest) (*attendanceModel.Attendance, error)
	AddAttendanceBulk(ctx context.Context, req *attendanceModel.CreateAttendanceRequest) (*attendanceModel.CreateAttendanceResponse, error)
	DeleteAttendanceByID(ctx context.Context, req *attendanceModel.DeleteAttendanceByIDRequest) (*attendanceModel.DeleteAttendanceByIDResponse, error)
	CheckInOutAttendance(ctx context.Context, req *attendanceModel.BulkCheckInOutAttendanceRequest) (*attendanceModel.BulkCheckInOutAttendanceResponse, error)
}

type attendanceService struct {
	dbProvider *attendanceDb.AttendanceRepository
	logger     *logger.Logger
}

func NewAttendanceService(dbProvider *attendanceDb.AttendanceRepository, logger *logger.Logger) AttendanceService {
	return &attendanceService{
		dbProvider: dbProvider,
		logger:     logger,
	}
}

func (s *attendanceService) ListAttendance(ctx context.Context, req *attendanceModel.ListAttendanceRequest) (*attendanceModel.ListAttendanceResponse, error) {
	funcName := "AddGuest"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if err := utils.ValidateUserAuthorization(currentUser, req.ProjectID, req.EventId, loggerZap); err != nil {
		return nil, err
	}

	if req.EventId == "" || req.ProjectID == "" {
		loggerZap.Error("err Invalid Argument missing event_id or project_id", nil)
		return nil, status.Error(codes.InvalidArgument, "Invalid Argument: missing event_id or project_id")
	}

	loggerZap.Info("Start ListProducts with req : ", req)
	loggerZap.Info("Start Decode Filter")

	decodeQuery, err := base64.RawStdEncoding.DecodeString(req.Query)
	if err != nil {
		loggerZap.Error("err DecodeString ", err)
		return nil, nil
	}

	loggerZap.Info("Success Decode Query")

	pagination := utils.SetPagination(req.Page, req.Limit)

	allowedColumns := map[string]bool{
		"created_at": true,
		"name":       true,
		"address":    true,
		"phone":      true,
		"email":      true,
	}

	allowedDirections := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	column := strings.ToLower(req.Sort)
	direction := strings.ToLower(req.Dir)

	if column != "" || direction != "" {

		if !allowedColumns[column] {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid Argument")
		}
		if !allowedDirections[direction] {
			return nil, status.Errorf(codes.InvalidArgument, "Invalid Argument")
		}
	}
	sort := &model.Sort{
		Column:    column,
		Direction: direction,
	}

	sqlBuilder := &db.QueryBuilder{
		CollectiveAnd: string(decodeQuery),
		Sort:          sort,
	}

	loggerZap.Info("Start ListAttendance")
	attendance, err := s.dbProvider.ListAttendance(ctx, req, pagination, sqlBuilder, sort)
	if err != nil {
		s.logger.Error("err ListGuests ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Start making response")

	result := &attendanceModel.ListAttendanceResponse{
		Error:      false,
		Code:       http.StatusOK,
		Message:    "Success",
		Data:       attendance,
		Pagination: pagination,
	}

	return result, nil
}

func (s *attendanceService) GetAttendanceByID(ctx context.Context, req *attendanceModel.GetAttendanceByIDRequest) (*attendanceModel.Attendance, error) {
	funcName := "GetAttendanceByID"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if err := utils.ValidateUserAuthorization(currentUser, req.ProjectID, req.EventId, loggerZap); err != nil {
		return nil, err
	}

	loggerZap.Info("Start GetAttendanceByID with req : ", req)

	attendance, err := s.dbProvider.GetAttendanceByID(ctx, req)
	if err != nil {
		s.logger.Error("err GetAttendanceByID ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	return attendance, nil
}

func (s *attendanceService) AddAttendanceBulk(ctx context.Context, req *attendanceModel.CreateAttendanceRequest) (*attendanceModel.CreateAttendanceResponse, error) {
	funcName := "AddAttendanceBulk"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if err := utils.ValidateUserAuthorization(currentUser, req.ProjectID, req.EventId, loggerZap); err != nil {
		return nil, err
	}

	loggerZap.Info("Start AddAttendanceBulk with req : ", req)

	listGuest, listGuestErr := s.dbProvider.GetListGuestByEventID(ctx, currentUser, req)
	if listGuestErr != nil {
		s.logger.Error("err GetListGuestByEventID ", listGuestErr)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	validGuestMap := make(map[string]*guestModel.Guest)
	for _, guest := range listGuest {
		guestIDStr := fmt.Sprint(guest.GuestID)
		validGuestMap[guestIDStr] = guest
	}

	validGuestIDs := make([]*string, 0)
	guestData := make(map[string]*guestModel.Guest)
	results := make([]*attendanceModel.GuestAttendanceResult, 0)

	// First pass: validate guest existence
	guestIDsToCheck := make([]int64, 0)
	guestIDMap := make(map[string]int64)

	for _, guestID := range req.GuestID {
		if guestID != nil {
			if guest, ok := validGuestMap[*guestID]; ok {
				guestIDsToCheck = append(guestIDsToCheck, guest.GuestID)
				guestIDMap[*guestID] = guest.GuestID
			} else {
				results = append(results, &attendanceModel.GuestAttendanceResult{
					GuestID: *guestID,
					Error:   true,
					Message: "Unauthorized: Guest ID not found in this project and event",
				})
			}
		}
	}

	// Check for existing attendance records
	var eventID, projectID int64
	{
		eventIDInt, _ := strconv.ParseInt(req.EventId, 0, 64)
		projectIDInt, _ := strconv.ParseInt(req.ProjectID, 0, 64)
		eventID = eventIDInt
		projectID = projectIDInt
	}

	existingGuestIDs := make(map[int64]bool)
	if len(guestIDsToCheck) > 0 {
		existing, existingErr := s.dbProvider.GetExistingAttendanceByGuestIDs(ctx, guestIDsToCheck, eventID, projectID)
		if existingErr != nil {
			s.logger.Error("err GetExistingAttendanceByGuestIDs ", existingErr)
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
		for _, id := range existing {
			existingGuestIDs[id] = true
		}
	}

	// Second pass: process valid guests that don't have existing attendance
	for _, guestID := range req.GuestID {
		if guestID != nil {
			if guest, ok := validGuestMap[*guestID]; ok {
				if existingGuestIDs[guest.GuestID] {
					results = append(results, &attendanceModel.GuestAttendanceResult{
						GuestID: *guestID,
						Error:   true,
						Message: "Guest already has attendance in this event",
					})
				} else {
					validGuestIDs = append(validGuestIDs, guestID)
					guestData[*guestID] = guest
					results = append(results, &attendanceModel.GuestAttendanceResult{
						GuestID: *guestID,
						Error:   false,
						Message: "Attendance created successfully",
					})
				}
			}
		}
	}

	if len(validGuestIDs) > 0 {
		validReq := &attendanceModel.CreateAttendanceRequest{
			EventId:   req.EventId,
			ProjectID: req.ProjectID,
			GuestID:   validGuestIDs,
			GuestData: guestData,
		}

		if err := s.dbProvider.CreateAttendance(ctx, validReq, currentUser); err != nil {
			s.logger.Error("err CreateAttendanceBulk ", err)
			for _, result := range results {
				if !result.Error {
					result.Error = true
					result.Message = "Failed to create attendance: Internal Server Error"
				}
			}
			return &attendanceModel.CreateAttendanceResponse{
				Error:   true,
				Code:    http.StatusInternalServerError,
				Message: "Failed to create attendance for valid guests",
				Data:    results,
			}, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	response := &attendanceModel.CreateAttendanceResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    results,
	}

	return response, nil
}

func (s *attendanceService) DeleteAttendanceByID(ctx context.Context, req *attendanceModel.DeleteAttendanceByIDRequest) (*attendanceModel.DeleteAttendanceByIDResponse, error) {
	funcName := "DeleteAttendanceByID"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if err := utils.ValidateUserAuthorization(currentUser, req.ProjectID, req.EventID, loggerZap); err != nil {
		return nil, err
	}

	loggerZap.Info("Start DeleteAttendanceByID with req : ", req)

	attendanceRecords, err := s.dbProvider.GetAttendancesByIDs(ctx, req)
	if err != nil {
		s.logger.Error("err GetAttendancesByIDs ", err)
		results := make([]*attendanceModel.AttendanceDeleteResult, 0)
		for _, attendanceID := range req.AttendanceIDs {
			results = append(results, &attendanceModel.AttendanceDeleteResult{
				AttendanceID: attendanceID,
				Error:        true,
				Message:      "Internal Server Error",
			})
		}
		return &attendanceModel.DeleteAttendanceByIDResponse{
			Error:   true,
			Code:    http.StatusInternalServerError,
			Message: "Failed to fetch attendance records",
			Data:    results,
		}, nil
	}

	foundAttendanceMap := make(map[string]bool)
	validAttendanceIDs := make([]string, 0)
	for _, attendance := range attendanceRecords {
		attendanceIDStr := fmt.Sprintf("%d", attendance.ID)
		foundAttendanceMap[attendanceIDStr] = true
		validAttendanceIDs = append(validAttendanceIDs, attendanceIDStr)
	}

	results := make([]*attendanceModel.AttendanceDeleteResult, 0)
	for _, attendanceID := range req.AttendanceIDs {
		if foundAttendanceMap[attendanceID] {
			results = append(results, &attendanceModel.AttendanceDeleteResult{
				AttendanceID: attendanceID,
				Error:        false,
				Message:      "Success",
			})
		} else {
			results = append(results, &attendanceModel.AttendanceDeleteResult{
				AttendanceID: attendanceID,
				Error:        true,
				Message:      "Unauthorized",
			})
		}
	}

	if len(validAttendanceIDs) > 0 {
		deleteErr := s.dbProvider.DeleteAttendanceByIDs(ctx, req)
		if deleteErr != nil {
			s.logger.Error("err DeleteAttendanceByIDs ", deleteErr)
			for _, result := range results {
				if !result.Error {
					result.Error = true
					result.Message = "Error when deleting"
				}
			}
			return &attendanceModel.DeleteAttendanceByIDResponse{
				Error:   true,
				Code:    http.StatusInternalServerError,
				Message: "Failed to delete some attendance records",
				Data:    results,
			}, nil
		}
	}

	hasError := false
	for _, result := range results {
		if result.Error {
			hasError = true
			break
		}
	}

	response := &attendanceModel.DeleteAttendanceByIDResponse{
		Error:   hasError,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    results,
	}

	return response, nil
}

func (s *attendanceService) CheckInOutAttendance(ctx context.Context, req *attendanceModel.BulkCheckInOutAttendanceRequest) (*attendanceModel.BulkCheckInOutAttendanceResponse, error) {
	funcName := "CheckInOutAttendance"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if err := utils.ValidateUserAuthorization(currentUser, req.ProjectID, req.EventID, loggerZap); err != nil {
		return nil, err
	}

	loggerZap.Info("Start CheckInOutAttendance with req : ", req)

	attendanceRecords, fetchErr := s.dbProvider.GetAttendanceRecordsByIDs(ctx, req.AttendanceIDs, req.ProjectID, req.EventID)
	if fetchErr != nil {
		s.logger.Error("err GetAttendanceRecordsByIDs ", fetchErr)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	attendanceMap := make(map[string]*attendanceModel.Attendance)
	for _, att := range attendanceRecords {
		attendanceMap[fmt.Sprintf("%d", att.ID)] = att
	}

	validAttendanceIDs := make([]string, 0)
	results := make([]*attendanceModel.CheckInOutResult, 0)

	for _, attendanceID := range req.AttendanceIDs {
		att, exists := attendanceMap[attendanceID]
		if !exists {
			results = append(results, &attendanceModel.CheckInOutResult{
				AttendanceID: attendanceID,
				Error:        true,
				Message:      "Unauthorized: Attendance ID not found in this project and event",
			})
			continue
		}

		isValid := false
		var errorMsg string

		switch req.Type {
		case constant.AttendanceTypeCheckIn:
			if att.CheckedInAt != nil {
				errorMsg = fmt.Sprintf("Already checked in at %s", att.CheckedInAt.Format("2006-01-02 15:04:05"))
			} else {
				isValid = true
			}
		case constant.AttendanceTypeCheckOut:
			if att.CheckedOutAt != nil {
				errorMsg = fmt.Sprintf("Already checked out at %s", att.CheckedOutAt.Format("2006-01-02 15:04:05"))
			} else {
				isValid = true
			}
		}

		if isValid {
			validAttendanceIDs = append(validAttendanceIDs, attendanceID)
			message := "Check-in successful"
			if req.Type == constant.AttendanceTypeCheckOut {
				message = "Check-out successful"
			}

			results = append(results, &attendanceModel.CheckInOutResult{
				AttendanceID: attendanceID,
				Error:        false,
				Message:      message,
			})
		} else {
			results = append(results, &attendanceModel.CheckInOutResult{
				AttendanceID: attendanceID,
				Error:        true,
				Message:      errorMsg,
			})
		}
	}

	if len(validAttendanceIDs) > 0 {

		checkInOutReq := &attendanceModel.BulkCheckInOutAttendanceRequest{
			EventID:       req.EventID,
			ProjectID:     req.ProjectID,
			AttendanceIDs: validAttendanceIDs,
			Type:          req.Type,
		}

		checkInOutErr := s.dbProvider.CheckInOutAttendance(ctx, checkInOutReq, currentUser)
		if checkInOutErr != nil {
			s.logger.Error("err CheckInOutAttendance ", checkInOutErr)

			for _, result := range results {
				if !result.Error {
					result.Error = true
					result.Message = "Failed to check-in/out: Internal Server Error"
				}
			}
			return &attendanceModel.BulkCheckInOutAttendanceResponse{
				Error:   true,
				Code:    http.StatusInternalServerError,
				Message: "Failed to check-in/out for valid attendances",
				Data:    results,
			}, nil
		}
	}

	response := &attendanceModel.BulkCheckInOutAttendanceResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    results,
	}

	return response, nil
}
