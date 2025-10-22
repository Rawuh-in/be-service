package event_service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	eventModel "rawuh-service/internal/event/model"
	eventDb "rawuh-service/internal/event/repository"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/model"
	"rawuh-service/internal/shared/redis"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type EventService interface {
	ListEvent(ctx context.Context, req *eventModel.ListEventRequest) (*eventModel.ListEventResponse, error)
	DetailEvent(ctx context.Context, req *eventModel.DetailEventRequest) (*eventModel.DetailEventResponse, error)
	DeleteEvent(ctx context.Context, req *eventModel.DeleteEventRequest) (*eventModel.DeleteEventResponse, error)
	AddEvent(ctx context.Context, req *eventModel.CreateEventRequest) error
	UpdateEvent(ctx context.Context, req *eventModel.UpdateEventRequest) error
}

type eventService struct {
	dbProvider *eventDb.EventRepository
	logger     *logrus.Logger
	redis      *redis.Redis
}

func NewEventService(dbProvider *eventDb.EventRepository, logger *logrus.Logger, redis *redis.Redis) EventService {
	return &eventService{
		dbProvider: dbProvider,
		logger:     logger,
		redis:      redis,
	}
}

func (s *eventService) ListEvent(ctx context.Context, req *eventModel.ListEventRequest) (*eventModel.ListEventResponse, error) {
	s.logger.Info("Start ListEvent with req : ", req)
	s.logger.Info("Start Decode Filter")

	if req.UserID == "" {
		s.logger.Error("Access denied with user_id  : ", req.UserID)
		return nil, status.Errorf(codes.PermissionDenied, "Permission Denied")
	}

	decodeQuery, err := base64.RawStdEncoding.DecodeString(req.Query)
	if err != nil {
		s.logger.Error("err DecodeString ", err)
		return nil, nil
	}

	s.logger.Info("Success Decode Query")

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

	s.logger.Info("Start ListEvent")
	guest, err := s.dbProvider.ListEvent(ctx, req.UserID, req.UserID, pagination, sqlBuilder, sort)
	if err != nil {
		s.logger.Error("err ListEvent ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Start making response")

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
	s.logger.Info("Start ListEvent with req : ", req)
	s.logger.Info("Start Decode Filter")

	if req.UserID == "" || req.EventsID == "" {
		s.logger.Error("Access denied with user_id  : ", req.UserID)
		return nil, status.Errorf(codes.PermissionDenied, "Permission Denied")
	}

	s.logger.Info("Start ListEvent")
	event, err := s.dbProvider.GetEventByID(ctx, req.EventsID, req.UserID)
	if err != nil {
		s.logger.Error("err ListEvent ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Start making response")

	result := &eventModel.DetailEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    event,
	}

	return result, nil
}

func (s *eventService) DeleteEvent(ctx context.Context, req *eventModel.DeleteEventRequest) (*eventModel.DeleteEventResponse, error) {
	s.logger.Info("Start ListEvent with req : ", req)
	s.logger.Info("Start Decode Filter")

	if req.UserID == "" || req.EventsID == "" {
		s.logger.Error("Access denied with user_id  : ", req.UserID)
		return nil, status.Errorf(codes.PermissionDenied, "Permission Denied")
	}

	s.logger.Info("Start ListEvent")
	err := s.dbProvider.DeleteEventByID(ctx, req.EventsID, req.UserID)
	if err != nil {
		s.logger.Error("err ListEvent ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}
	result := &eventModel.DeleteEventResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
	}

	return result, nil
}

func (s *eventService) AddEvent(ctx context.Context, req *eventModel.CreateEventRequest) error {
	remarkLength, _ := strconv.Atoi(utils.GetEnv("REMARK_LENGTH", "500"))
	nameLength, _ := strconv.Atoi(utils.GetEnv("PRODUCT_NAME_LENGTH", "255"))

	s.logger.Info("Start Validation for req ", req)

	if utils.IsEmptyString(req.EventName) || strings.TrimSpace(req.UserID) == "" {
		return status.Errorf(codes.Aborted, "guest name is empty")
	}
	if len(req.EventName) > nameLength {
		return status.Errorf(codes.Aborted, "guest name maximum characters is %d", nameLength)
	}

	if !utils.IsValidProductName(req.EventName) || strings.ContainsAny(req.EventName, "%$#@!*&^<>\"") || !utils.IsValidProductName(req.Description) {
		return status.Errorf(codes.Aborted, "characters not allowed in guest name")
	}

	if strings.TrimSpace(req.Description) != "" {

		if len(req.Description) > remarkLength {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprintf("%s maximum characters is %d", req.Description, remarkLength))
		}
		if !utils.IsValidCharacter(req.Description) {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprint("characters not allowed in field Address", req.Description))
		}
	}

	s.logger.Info("Start AddEvent with data ", req)

	optionData, _ := json.Marshal(req.Options)
	// optionData, _ := protojson.MarshalOptions{UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(req.Options)

	req.Options = string(optionData)

	err := s.dbProvider.CreateEvent(ctx, req)
	if err != nil {
		s.logger.Error("err AddEvent ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Success AddEvent")

	return nil
}

func (s *eventService) UpdateEvent(ctx context.Context, req *eventModel.UpdateEventRequest) error {
	remarkLength, _ := strconv.Atoi(utils.GetEnv("REMARK_LENGTH", "500"))
	nameLength, _ := strconv.Atoi(utils.GetEnv("PRODUCT_NAME_LENGTH", "255"))

	s.logger.Info("Start Validation for req ", req)

	if utils.IsEmptyString(req.EventName) || strings.TrimSpace(req.UserID) == "" {
		return status.Errorf(codes.Aborted, "guest name is empty")
	}
	if len(req.EventName) > nameLength {
		return status.Errorf(codes.Aborted, "guest name maximum characters is %d", nameLength)
	}

	if !utils.IsValidProductName(req.EventName) || strings.ContainsAny(req.EventName, "%$#@!*&^<>\"") || !utils.IsValidProductName(req.Description) {
		return status.Errorf(codes.Aborted, "characters not allowed in guest name")
	}

	if strings.TrimSpace(req.Description) != "" {

		if len(req.Description) > remarkLength {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprintf("%s maximum characters is %d", req.Description, remarkLength))
		}
		if !utils.IsValidCharacter(req.Description) {
			return status.Errorf(codes.Aborted, "%s", fmt.Sprint("characters not allowed in field Address", req.Description))
		}
	}

	s.logger.Info("Start AddEvent with data ", req)

	optionData, _ := json.Marshal(req.Options)
	// optionData, _ := protojson.MarshalOptions{UseEnumNumbers: true, EmitUnpopulated: true}.Marshal(req.Options)

	req.Options = string(optionData)

	err := s.dbProvider.UpdateEvent(ctx, req)
	if err != nil {
		s.logger.Error("err AddEvent ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Success AddEvent")

	return nil
}
