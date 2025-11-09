package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	guestModel "rawuh-service/internal/guest/model"
	"rawuh-service/internal/shared/constant"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/logger"
	"rawuh-service/internal/shared/middleware"
	"rawuh-service/internal/shared/model"
	"strconv"
	"strings"

	guestDb "rawuh-service/internal/guest/repository"
	db "rawuh-service/internal/shared/db"

	"go.elastic.co/apm/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type GuestService interface {
	AddGuest(ctx context.Context, p *guestModel.CreateGuestRequest) error
	UpdateGuestByID(ctx context.Context, p *guestModel.UpdateGuestRequest) error
	GetGuestByID(ctx context.Context, req *guestModel.GetGuestByIDRequest) (*guestModel.GetGuestByIDResponse, error)
	DeleteGuestByID(ctx context.Context, req *guestModel.DeleteGuestByIDRequest) error
	ListGuests(ctx context.Context, req *guestModel.ListGuestRequest) (*guestModel.ListGuestResponse, error)
}

type guestService struct {
	dbProvider *guestDb.GuestRepository
	logger     *logger.Logger
	// redis      *redis.Redis
}

func NewGuestService(dbProvider *guestDb.GuestRepository, logger *logger.Logger) GuestService {
	return &guestService{
		dbProvider: dbProvider,
		logger:     logger,
		// redis:      redis,
	}
}

func (s *guestService) AddGuest(ctx context.Context, req *guestModel.CreateGuestRequest) error {
	funcName := "AddGuest"
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
		if req.ProjectID != fmt.Sprintf("%d", currentUser.ProjectID) || req.EventId != fmt.Sprintf("%d", currentUser.EventID) {
			loggerZap.Error("err GetMeFromMD unauthorized user", nil)
			return status.Error(codes.PermissionDenied, "Permission Denied")
		}
	default:
		loggerZap.Error("err GetMeFromMD unauthorized user type", nil)
		return status.Error(codes.PermissionDenied, "Permission Denied")
	}

	remarkLength, _ := strconv.Atoi(utils.GetEnv("GUEST_REMARK_LENGTH", "500"))
	nameLength, _ := strconv.Atoi(utils.GetEnv("GUEST_NAME_LENGTH", "255"))

	loggerZap.Info("Start Validation for req ", req)

	if utils.IsEmptyString(req.Name) {
		return status.Errorf(codes.Aborted, "guest name is empty")
	}
	if len(req.Name) > nameLength {
		return status.Errorf(codes.Aborted, "guest name maximum characters is %d", nameLength)
	}

	if !utils.IsValidProductName(req.Name) {
		return status.Errorf(codes.Aborted, "characters not allowed in guest name")
	}

	if strings.TrimSpace(req.Address) != "" {

		if len(req.Address) > remarkLength {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprintf("%s maximum characters is %d", req.Address, remarkLength))
		}
		if !utils.IsValidCharacter(req.Address) {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprint("characters not allowed in field Address", req.Address))
		}
	}

	loggerZap.Info("Start CreateGuest with data ", req)

	if req.EventData != "" {
		var optionStr map[string]interface{}
		if err := json.Unmarshal([]byte(req.EventData), &optionStr); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}

		utils.SanitizeJSON(optionStr)
		optionData, _ := json.Marshal(optionStr)
		req.EventData = string(optionData)
	} else {
		req.EventData = "{}"
	}

	if req.GuestData != "" {
		var optionStr map[string]interface{}
		if err := json.Unmarshal([]byte(req.GuestData), &optionStr); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}

		utils.SanitizeJSON(optionStr)
		optionData, _ := json.Marshal(optionStr)
		req.GuestData = string(optionData)
	} else {
		req.GuestData = "{}"
	}

	err := s.dbProvider.CreateGuest(ctx, req, currentUser)
	if err != nil {
		loggerZap.Error("err CreateGuest ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Success CreateGuest")

	return nil
}

func (s *guestService) UpdateGuestByID(ctx context.Context, req *guestModel.UpdateGuestRequest) error {
	funcName := "UpdateGuestByID"
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
		if req.ProjectID != fmt.Sprintf("%d", currentUser.ProjectID) || req.EventId != fmt.Sprintf("%d", currentUser.EventID) {
			loggerZap.Error("err GetMeFromMD unauthorized user", nil)
			return status.Error(codes.PermissionDenied, "Permission Denied")
		}
	default:
		loggerZap.Error("err GetMeFromMD unauthorized user type", nil)
		return status.Error(codes.PermissionDenied, "Permission Denied")
	}

	remarkLength, _ := strconv.Atoi(utils.GetEnv("GUEST_REMARK_LENGTH", "500"))
	nameLength, _ := strconv.Atoi(utils.GetEnv("GUEST_NAME_LENGTH", "255"))

	loggerZap.Info("Start Validation for req ", req)

	if utils.IsEmptyString(req.Name) {
		return status.Errorf(codes.Aborted, "guest name is empty")
	}
	if len(req.Name) > nameLength {
		return status.Errorf(codes.Aborted, "guest name maximum characters is %d", nameLength)
	}

	if !utils.IsValidProductName(req.Name) {
		return status.Errorf(codes.Aborted, "characters not allowed in guest name")
	}

	if req.EventId == "" {
		return status.Errorf(codes.Aborted, "invalid event id")
	}

	if strings.TrimSpace(req.Address) != "" {

		if len(req.Address) > remarkLength {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprintf("%s maximum characters is %d", req.Address, remarkLength))
		}
		if !utils.IsValidCharacter(req.Address) {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprint("characters not allowed in field Address", req.Address))
		}
	}

	if req.EventData != "" {
		var optionStr map[string]interface{}
		if err := json.Unmarshal([]byte(req.EventData), &optionStr); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}

		utils.SanitizeJSON(optionStr)
		optionData, _ := json.Marshal(optionStr)
		req.EventData = string(optionData)
	} else {
		req.EventData = "{}"
	}

	if req.GuestData != "" {
		var optionStr map[string]interface{}
		if err := json.Unmarshal([]byte(req.GuestData), &optionStr); err != nil {
			return fmt.Errorf("invalid JSON format: %w", err)
		}

		utils.SanitizeJSON(optionStr)
		optionData, _ := json.Marshal(optionStr)
		req.GuestData = string(optionData)
	} else {
		req.GuestData = "{}"
	}

	loggerZap.Info("Start UpdateGuest with data ", req)

	err := s.dbProvider.UpdateGuest(ctx, req, currentUser)
	if err != nil {
		loggerZap.Error("err UpdateGuest ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Success UpdateGuest")

	return nil
}

func (s *guestService) ListGuests(ctx context.Context, req *guestModel.ListGuestRequest) (*guestModel.ListGuestResponse, error) {
	funcName := "ListGuests"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	switch currentUser.UserType {
	case constant.UserTypeSystemAdmin:
		// system admin can access all projects
	case constant.UserTypeProjectUser:
		if req.ProjectID != fmt.Sprintf("%d", currentUser.ProjectID) || req.EventId != fmt.Sprintf("%d", currentUser.EventID) {
			loggerZap.Error("err GetMeFromMD unauthorized user", nil)
			return nil, status.Error(codes.PermissionDenied, "Permission Denied")
		}
	default:
		loggerZap.Error("err GetMeFromMD unauthorized user type", nil)
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
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

	loggerZap.Info("Start ListGuests")
	guest, err := s.dbProvider.ListGuests(ctx, req, pagination, sqlBuilder, sort)
	if err != nil {
		s.logger.Error("err ListGuests ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Start making response")

	result := &guestModel.ListGuestResponse{
		Error:      false,
		Code:       http.StatusOK,
		Message:    "Success",
		Data:       guest,
		Pagination: pagination,
	}

	return result, nil

}
func (s *guestService) GetGuestByID(ctx context.Context, req *guestModel.GetGuestByIDRequest) (*guestModel.GetGuestByIDResponse, error) {
	funcName := "GetGuestByID"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.GuestID == "" {
		loggerZap.Error("err Invalid event id : ", nil)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Event Id")
	}

	switch currentUser.UserType {
	case constant.UserTypeSystemAdmin:
		// system admin can access all projects
	case constant.UserTypeProjectUser:
		if req.ProjectID != fmt.Sprintf("%d", currentUser.ProjectID) || req.EventId != fmt.Sprintf("%d", currentUser.EventID) {
			loggerZap.Error("err GetMeFromMD unauthorized user", nil)
			return nil, status.Error(codes.PermissionDenied, "Permission Denied")
		}
	default:
		loggerZap.Error("err GetMeFromMD unauthorized user type", nil)
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	loggerZap.Info("Start GetGuestByID with req : ", req)

	loggerZap.Info("Start GetGuestByID")
	guest, err := s.dbProvider.GetGuestByID(ctx, req)
	if err != nil {
		loggerZap.Error("err GetGuestByID ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	if guest == nil || guest.ProjectID == 0 {
		loggerZap.Info("GetGuestByID not found", nil)
		return nil, status.Errorf(codes.NotFound, "guest not found")
	}

	loggerZap.Info("Start making response")

	result := &guestModel.GetGuestByIDResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    guest,
	}

	return result, nil

}
func (s *guestService) DeleteGuestByID(ctx context.Context, req *guestModel.DeleteGuestByIDRequest) error {
	funcName := "ListProjects"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	if req.GuestID == "" {
		loggerZap.Error("err Invalid event id : ", nil)
		return status.Errorf(codes.InvalidArgument, "Invalid Event Id")
	}

	switch currentUser.UserType {
	case constant.UserTypeSystemAdmin:
		// system admin can access all projects
	case constant.UserTypeProjectUser:
		if req.ProjectID != fmt.Sprintf("%d", currentUser.ProjectID) || req.EventId != fmt.Sprintf("%d", currentUser.EventID) {
			loggerZap.Error("err GetMeFromMD unauthorized user", nil)
			return status.Error(codes.PermissionDenied, "Permission Denied")
		}
	default:
		loggerZap.Error("err GetMeFromMD unauthorized user type", nil)
		return status.Error(codes.PermissionDenied, "Permission Denied")
	}

	loggerZap.Info("Start DeleteGuestByID with req : ", req)

	loggerZap.Info("Start DeleteGuestByID")
	err := s.dbProvider.DeleteGuestByID(ctx, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			loggerZap.Warn("guest not found", err)
			return status.Error(codes.NotFound, "Guest not found")
		}

		s.logger.Error("err DeleteGuestByID ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Start making response")

	return nil

}
