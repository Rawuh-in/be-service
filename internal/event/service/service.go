package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	eventModel "rawuh-service/internal/event/model"
	eventDb "rawuh-service/internal/event/repository"
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
	"gorm.io/gorm"
)

type EventService interface {
	ListEvent(ctx context.Context, req *eventModel.ListEventRequest) (*eventModel.ListEventResponse, error)
	DetailEvent(ctx context.Context, req *eventModel.DetailEventRequest) (*eventModel.DetailEventResponse, error)
	DeleteEvent(ctx context.Context, req *eventModel.DeleteEventRequest) error
	AddEvent(ctx context.Context, req *eventModel.CreateEventRequest) error
	UpdateEvent(ctx context.Context, req *eventModel.UpdateEventRequest) error
}

type eventService struct {
	dbProvider *eventDb.EventRepository
	logger     *logger.Logger
	// redis      *redis.Redis
}

func NewEventService(dbProvider *eventDb.EventRepository, logger *logger.Logger) EventService {
	return &eventService{
		dbProvider: dbProvider,
		logger:     logger,
		// redis:      redis,
	}
}

func (s *eventService) ListEvent(ctx context.Context, req *eventModel.ListEventRequest) (*eventModel.ListEventResponse, error) {
	funcName := "ListEvent"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	loggerZap.Info("Start ListEvent with req : ", req)
	loggerZap.Info("Start Decode Filter")

	decodeQuery, err := base64.RawStdEncoding.DecodeString(req.Query)
	if err != nil {
		loggerZap.Error("err DecodeString ", err)
		return nil, nil
	}

	loggerZap.Info("Success Decode Query")

	pagination := utils.SetPagination(req.Page, req.Limit)

	allowedColumns := map[string]bool{
		"event_name": true,
		"event_id":   true,
		"created_at": true,
		"updated_at": true,
		"start_date": true,
		"end_date":   true,
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

	loggerZap.Info("Start ListEvent")
	guest, err := s.dbProvider.ListEvent(ctx, currentUser, pagination, sqlBuilder, sort)
	if err != nil {
		loggerZap.Error("err ListEvent ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Start making response")

	result := &eventModel.ListEventResponse{
		Error:      false,
		Code:       http.StatusOK,
		Message:    "Success",
		Data:       guest,
		Pagination: pagination,
	}

	return result, nil
}

func (s *eventService) DetailEvent(ctx context.Context, req *eventModel.DetailEventRequest) (*eventModel.DetailEventResponse, error) {
	funcName := "ListProjects"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	loggerZap.Info("Start ListEvent with req : ", req)

	if req.EventsID == "" {
		loggerZap.Error("Empty Event Id  : ", nil)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Argument")
	}

	switch currentUser.UserType {
	case constant.UserTypeSystemAdmin:
		// system admin can access all projects
	case constant.UserTypeProjectUser:
		if req.ProjectID != fmt.Sprintf("%d", currentUser.ProjectID) {
			loggerZap.Error("err GetMeFromMD unauthorized user", nil)
			return nil, status.Error(codes.PermissionDenied, "Permission Denied")
		}
	default:
		loggerZap.Error("err GetMeFromMD unauthorized user type", nil)
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	loggerZap.Info("Start ListEvent")
	event, err := s.dbProvider.GetEventByID(ctx, req.EventsID, currentUser)
	if err != nil {
		loggerZap.Error("err ListEvent ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	if event == nil || event.EventID == 0 {
		loggerZap.Info("event not found", nil)
		return nil, status.Errorf(codes.NotFound, "event not found")
	}

	loggerZap.Info("Start making response")

	result := &eventModel.DetailEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    event,
	}

	return result, nil
}

func (s *eventService) DeleteEvent(ctx context.Context, req *eventModel.DeleteEventRequest) error {
	funcName := "DeleteEvent"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	loggerZap.Info("Start ListEvent with req : ", req)
	loggerZap.Info("Start Decode Filter")

	if req.EventsID == "" {
		loggerZap.Error("Access denied with user_id  : ", nil)
		return status.Errorf(codes.PermissionDenied, "Permission Denied")
	}

	switch currentUser.UserType {
	case constant.UserTypeSystemAdmin:
		// system admin can access all projects
	case constant.UserTypeProjectUser:
		if req.ProjectID != fmt.Sprintf("%d", currentUser.ProjectID) {
			loggerZap.Error("err GetMeFromMD unauthorized user", nil)
			return status.Error(codes.PermissionDenied, "Permission Denied")
		}
	default:
		loggerZap.Error("err GetMeFromMD unauthorized user type", nil)
		return status.Error(codes.PermissionDenied, "Permission Denied")
	}

	loggerZap.Info("Start ListEvent")
	err := s.dbProvider.DeleteEventByID(ctx, req.EventsID, currentUser)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			loggerZap.Warn("event not found", err)
			return status.Error(codes.NotFound, "Event not found")
		}

		s.logger.Error("err DeleteEventByID ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	return nil
}

func (s *eventService) AddEvent(ctx context.Context, req *eventModel.CreateEventRequest) error {
	funcName := "AddEvent"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	switch currentUser.UserType {
	case constant.UserTypeSystemAdmin:
		// system admin can access all projects
	case constant.UserTypeProjectUser:
		if req.ProjectID != fmt.Sprintf("%d", currentUser.ProjectID) {
			loggerZap.Error("err GetMeFromMD unauthorized user", nil)
			return status.Error(codes.PermissionDenied, "Permission Denied")
		}
	default:
		loggerZap.Error("err GetMeFromMD unauthorized user type", nil)
		return status.Error(codes.PermissionDenied, "Permission Denied")
	}

	remarkLength, _ := strconv.Atoi(utils.GetEnv("EVENT_REMARK_LENGTH", "500"))
	nameLength, _ := strconv.Atoi(utils.GetEnv("EVENT_NAME_LENGTH", "255"))

	loggerZap.Info("Start Validation for req ", req)

	if utils.IsEmptyString(req.EventName) {
		return status.Errorf(codes.Aborted, "event name is empty")
	}
	if len(req.EventName) > nameLength {
		return status.Errorf(codes.Aborted, "event name maximum characters is %d", nameLength)
	}

	if !utils.IsValidProductName(req.EventName) || strings.ContainsAny(req.EventName, "%$#@!*&^<>\"") || !utils.IsValidProductName(req.Description) {
		return status.Errorf(codes.Aborted, "characters not allowed in event name")
	}

	if strings.TrimSpace(req.Description) != "" {

		if len(req.Description) > remarkLength {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprintf("%s maximum characters is %d", req.Description, remarkLength))
		}
		if !utils.IsValidCharacter(req.Description) {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprint("characters not allowed in field Address", req.Description))
		}
	}

	loggerZap.Info("Start AddEvent with data ", req)

	if req.EventOptions != "" {
		var optionStr map[string]interface{}
		if err := json.Unmarshal([]byte(req.EventOptions), &optionStr); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}

		utils.SanitizeJSON(optionStr)
		optionData, _ := json.Marshal(optionStr)
		req.EventOptions = string(optionData)
	} else {
		req.EventOptions = "{}"
	}
	if req.GuestOptions != "" {
		var optionStr map[string]interface{}
		if err := json.Unmarshal([]byte(req.GuestOptions), &optionStr); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}

		utils.SanitizeJSON(optionStr)
		optionData, _ := json.Marshal(optionStr)
		req.GuestOptions = string(optionData)
	} else {
		req.GuestOptions = "{}"
	}

	err := s.dbProvider.CreateEvent(ctx, req, currentUser)
	if err != nil {
		loggerZap.Error("err AddEvent ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Success AddEvent")

	return nil
}

func (s *eventService) UpdateEvent(ctx context.Context, req *eventModel.UpdateEventRequest) error {
	funcName := "UpdateEvent"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	switch currentUser.UserType {
	case constant.UserTypeSystemAdmin:
		// system admin can access all projects
	case constant.UserTypeProjectUser:
		if req.ProjectID != fmt.Sprintf("%d", currentUser.ProjectID) {
			loggerZap.Error("err GetMeFromMD unauthorized user", nil)
			return status.Error(codes.PermissionDenied, "Permission Denied")
		}
	default:
		loggerZap.Error("err GetMeFromMD unauthorized user type", nil)
		return status.Error(codes.PermissionDenied, "Permission Denied")
	}

	remarkLength, _ := strconv.Atoi(utils.GetEnv("EVENT_REMARK_LENGTH", "500"))
	nameLength, _ := strconv.Atoi(utils.GetEnv("EVENT_NAME_LENGTH", "255"))

	loggerZap.Info("Start Validation for req ", req)

	if utils.IsEmptyString(req.EventName) {
		return status.Errorf(codes.Aborted, "event name is empty")
	}
	if len(req.EventName) > nameLength {
		return status.Errorf(codes.Aborted, "event name maximum characters is %d", nameLength)
	}

	if !utils.IsValidProductName(req.EventName) || strings.ContainsAny(req.EventName, "%$#@!*&^<>\"") || !utils.IsValidProductName(req.Description) {
		return status.Errorf(codes.Aborted, "characters not allowed in event name")
	}

	if strings.TrimSpace(req.Description) != "" {

		if len(req.Description) > remarkLength {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprintf("%s maximum characters is %d", req.Description, remarkLength))
		}
		if !utils.IsValidCharacter(req.Description) {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprint("characters not allowed in field Address", req.Description))
		}
	}
	loggerZap.Info("Start UpdateEvent with data ", req)

	if req.EventOptions != "" {
		var optionStr map[string]interface{}
		if err := json.Unmarshal([]byte(req.EventOptions), &optionStr); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}

		utils.SanitizeJSON(optionStr)
		optionData, _ := json.Marshal(optionStr)
		req.EventOptions = string(optionData)
	} else {
		req.EventOptions = "{}"
	}
	if req.GuestOptions != "" {
		var optionStr map[string]interface{}
		if err := json.Unmarshal([]byte(req.GuestOptions), &optionStr); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}

		utils.SanitizeJSON(optionStr)
		optionData, _ := json.Marshal(optionStr)
		req.GuestOptions = string(optionData)
	} else {
		req.GuestOptions = "{}"
	}

	err := s.dbProvider.UpdateEvent(ctx, req, currentUser)
	if err != nil {
		loggerZap.Error("err UpdateEvent ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Success UpdateEvent")

	return nil
}
