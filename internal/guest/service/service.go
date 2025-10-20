package guest_service

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	guest_model "rawuh-service/internal/guest/model"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/model"
	paginationModel "rawuh-service/internal/shared/model"
	"strconv"
	"strings"

	guest_db "rawuh-service/internal/guest/repository"
	db "rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/redis"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GuestService interface {
	AddGuest(ctx context.Context, p *guest_model.CreateGuestRequest) error
	ListGuests(ctx context.Context, req *guest_model.ListGuestRequest) (*guest_model.ListGuestResponse, error)
}

type guestService struct {
	dbProvider *guest_db.GuestRepository
	logger     *logrus.Logger
	redis      *redis.Redis
}

func NewGuestService(dbProvider *guest_db.GuestRepository, logger *logrus.Logger, redis *redis.Redis) GuestService {
	return &guestService{
		dbProvider: dbProvider,
		logger:     logger,
		redis:      redis,
	}
}

func setPagination(page int32, limit int32) *paginationModel.PaginationResponse {
	res := &paginationModel.PaginationResponse{
		Limit: 10,
		Page:  1,
	}

	if limit == 0 && page == 0 {
		res.Limit = -1
		res.Page = -1
		return res
	} else {
		res.Limit = limit
		res.Page = page
	}

	if res.Page == 0 {
		res.Page = 1
	}

	switch {
	case res.Limit > 100:
		res.Limit = 100
	case res.Limit <= 0:
		res.Limit = 10
	}

	return res
}

func (s *guestService) AddGuest(ctx context.Context, req *guest_model.CreateGuestRequest) error {
	remarkLength, _ := strconv.Atoi(utils.GetEnv("REMARK_LENGTH", "500"))
	nameLength, _ := strconv.Atoi(utils.GetEnv("PRODUCT_NAME_LENGTH", "255"))

	s.logger.Info("Start Validation for req ", req)

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
			return status.Errorf(codes.Aborted, fmt.Sprintf("%s maximum characters is %d", req.Address, remarkLength))
		}
		if !utils.IsValidCharacter(req.Address) {
			return status.Errorf(codes.Aborted, fmt.Sprintf("characters not allowed in field Address", req.Address))
		}
	}

	s.logger.Info("Start CreateGuest with data ", req)

	err := s.dbProvider.CreateGuest(ctx, req)
	if err != nil {
		s.logger.Error("err CreateGuest ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Success CreateGuest")

	return nil
}

func (s *guestService) ListGuests(ctx context.Context, req *guest_model.ListGuestRequest) (*guest_model.ListGuestResponse, error) {

	s.logger.Info("Start ListProducts with req : ", req)
	s.logger.Info("Start Decode Filter")

	if req.EventId == "" {
		s.logger.Error("err Invalid event id : ", req.EventId)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Event Id")
	}

	decodeQuery, err := base64.RawStdEncoding.DecodeString(req.Query)
	if err != nil {
		s.logger.Error("err DecodeString ", err)
		return nil, nil
	}

	s.logger.Info("Success Decode Query")

	pagination := setPagination(req.Page, req.Limit)

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

	s.logger.Info("Start ListGuests")
	guest, err := s.dbProvider.ListGuests(ctx, req.EventId, pagination, sqlBuilder, sort)
	if err != nil {
		s.logger.Error("err ListGuests ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Start making response")

	guests := &guest_model.ListGuestResponse{
		Error:      false,
		Code:       http.StatusOK,
		Message:    "Success",
		Data:       guest,
		Pagination: pagination,
	}

	return guests, nil

}
