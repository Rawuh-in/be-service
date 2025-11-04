package service

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	authModel "rawuh-service/internal/auth/model"
	repoAuth "rawuh-service/internal/auth/repository"
	"rawuh-service/internal/shared/constant"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/logger"
	"rawuh-service/internal/shared/middleware"
	"rawuh-service/internal/shared/model"
	"rawuh-service/internal/shared/redis"
	userModel "rawuh-service/internal/user/model"
	"strconv"
	"strings"

	db "rawuh-service/internal/shared/db"
	userDb "rawuh-service/internal/user/repository"

	"go.elastic.co/apm/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserService interface {
	AddUser(ctx context.Context, p *userModel.CreateUserRequest) error
	UpdateUserByID(ctx context.Context, p *userModel.UpdateUserRequest) error
	GetUserByID(ctx context.Context, req *userModel.GetUserByIDRequest) (*userModel.GetUserByIDResponse, error)
	DeleteUserByID(ctx context.Context, req *userModel.DeleteUserByIDRequest) error
	ListUsers(ctx context.Context, req *userModel.ListUserRequest) (*userModel.ListUserResponse, error)
}

type userService struct {
	dbProvider *userDb.UserRepository
	logger     *logger.Logger
	authRepo   *repoAuth.AuthRepository
	redis      *redis.Redis
}

func NewUserService(dbProvider *userDb.UserRepository, authRepo *repoAuth.AuthRepository, rdb *redis.Redis, logger *logger.Logger) UserService {
	return &userService{
		dbProvider: dbProvider,
		logger:     logger,
		authRepo:   authRepo,
		redis:      rdb,
	}
}

func (s *userService) AddUser(ctx context.Context, req *userModel.CreateUserRequest) error {
	funcName := "AddUser"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	loggerZap.Info("Success GetMeFromMD ", currentUser)

	if currentUser.UserType != constant.UserTypeSystemAdmin {
		loggerZap.Error("err CreateProject unauthorized user", nil)
		return status.Error(codes.PermissionDenied, "Permission Denied")
	}

	nameLength, _ := strconv.Atoi(utils.GetEnv("USER_NAME_LENGTH", "255"))

	loggerZap.Info("Start Validation for req ", req)

	if req.Name == "" || req.Username == "" || req.UserType == "" || req.ProjectID == "" {
		return status.Errorf(codes.Aborted, "invalid argument, missing required fields")
	}

	if req.Password == "" {
		return status.Errorf(codes.Aborted, "password required")
	}

	if len(req.Name) > nameLength {
		return status.Errorf(codes.Aborted, "user name maximum characters is %d", nameLength)
	}

	if !utils.IsValidProductName(req.Name) || !utils.IsValidProductName(req.Username) {
		return status.Errorf(codes.Aborted, "characters not allowed in user name")
	}

	loggerZap.Info("Start CreateUser with data ", req)

	found, err := s.dbProvider.CheckUsernameExist(ctx, req.Username)
	if err != nil {
		loggerZap.Error("err CheckUsernameExist ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}
	if found {
		return status.Errorf(codes.AlreadyExists, "username already exists")
	}

	// create user and get created user id
	createdID, err := s.dbProvider.CreateUser(ctx, req, currentUser)
	if err != nil {
		loggerZap.Error("err CreateUser ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	// // if password provided, create auth row
	// if req.Password != "" {
	// }
	encrypted, err := utils.EncryptAES(req.Password)
	if err != nil {
		loggerZap.Error("err EncryptAES", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	authRow := &authModel.Auth{
		UserID:    createdID,
		Username:  req.Username,
		Password:  encrypted,
		ProjectID: func() int64 { pid, _ := strconv.ParseInt(req.ProjectID, 0, 64); return pid }(),
	}

	if err := s.authRepo.CreateAuth(ctx, authRow); err != nil {
		loggerZap.Error("err CreateAuth", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Success CreateUser")

	return nil
}

func (s *userService) UpdateUserByID(ctx context.Context, req *userModel.UpdateUserRequest) error {
	funcName := "UpdateUserByID"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	loggerZap.Info("Success GetMeFromMD ", currentUser)

	if currentUser.UserType != constant.UserTypeSystemAdmin {
		loggerZap.Error("err CreateProject unauthorized user", nil)
		return status.Error(codes.PermissionDenied, "Permission Denied")
	}

	nameLength, _ := strconv.Atoi(utils.GetEnv("USER_NAME_LENGTH", "255"))

	loggerZap.Info("Start Validation for req ", req)

	if utils.IsEmptyString(req.Name) {
		return status.Errorf(codes.Aborted, "user name is empty")
	}
	if len(req.Name) > nameLength {
		return status.Errorf(codes.Aborted, "user name maximum characters is %d", nameLength)
	}

	if !utils.IsValidProductName(req.Name) {
		return status.Errorf(codes.Aborted, "characters not allowed in user name")
	}

	if req.UserID == "" {
		return status.Errorf(codes.Aborted, "invalid user id")
	}

	loggerZap.Info("Start UpdateUser with data ", req)

	err := s.dbProvider.UpdateUser(ctx, req, currentUser)
	if err != nil {
		loggerZap.Error("err UpdateUser ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Success UpdateUser")

	return nil
}

func (s *userService) ListUsers(ctx context.Context, req *userModel.ListUserRequest) (*userModel.ListUserResponse, error) {
	funcName := "ListUsers"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)

	loggerZap.Info("Start ListUsers with req : ", req)
	loggerZap.Info("Start Decode Filter")

	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	loggerZap.Info("Success GetMeFromMD ", currentUser)

	if currentUser.UserType != constant.UserTypeSystemAdmin {
		loggerZap.Error("err CreateProject unauthorized user", nil)
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

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
		"username":   true,
		"user_type":  true,
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

	loggerZap.Info("Start ListUsers")
	users, err := s.dbProvider.ListUsers(ctx, req.EventId, req.ProjectID, pagination, sqlBuilder, sort)
	if err != nil {
		s.logger.Error("err ListUsers ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Start making response")

	result := &userModel.ListUserResponse{
		Error:      false,
		Code:       http.StatusOK,
		Message:    "Success",
		Data:       users,
		Pagination: pagination,
	}

	return result, nil

}

func (s *userService) GetUserByID(ctx context.Context, req *userModel.GetUserByIDRequest) (*userModel.GetUserByIDResponse, error) {
	funcName := "GetUserByID"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)

	loggerZap.Info("Start GetUserByID with req : ", req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return nil, status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	loggerZap.Info("Success GetMeFromMD ", currentUser)

	if currentUser.UserType != constant.UserTypeSystemAdmin {
		loggerZap.Error("err CreateProject unauthorized user", nil)
		return nil, status.Error(codes.PermissionDenied, "Permission Denied")
	}

	if req.UserID == "" {
		loggerZap.Error("err Invalid event id : ", nil)
		return nil, status.Errorf(codes.InvalidArgument, "Invalid User Id")
	}

	loggerZap.Info("Start GetUserByID")
	user, err := s.dbProvider.GetUserByID(ctx, req.UserID)
	if err != nil {
		loggerZap.Error("err GetUserByID ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	if user == nil {
		loggerZap.Info("GetUserByID not found", nil)
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	loggerZap.Info("Start making response")

	result := &userModel.GetUserByIDResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    user,
	}

	return result, nil

}
func (s *userService) DeleteUserByID(ctx context.Context, req *userModel.DeleteUserByIDRequest) error {
	funcName := "DeleteUserByID"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	currentUser, ok := middleware.GetAuthClaimsFromContext(ctx)
	if !ok {
		loggerZap.Error("err GetMeFromMD no auth claims", nil)
		return status.Error(codes.Unauthenticated, "Unauthenticated")
	}

	loggerZap.Info("Success GetMeFromMD ", currentUser)

	if currentUser.UserType != constant.UserTypeSystemAdmin {
		loggerZap.Error("err CreateProject unauthorized user", nil)
		return status.Error(codes.PermissionDenied, "Permission Denied")
	}

	loggerZap.Info("Start DeleteUserByID with req : ", req)

	if req.UserID == "" {
		loggerZap.Error("err Invalid user id : ", nil)
		return status.Errorf(codes.InvalidArgument, "Invalid User Id")
	}

	loggerZap.Info("Start DeleteUserByID")
	err := s.dbProvider.DeleteUserByID(ctx, req.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			loggerZap.Warn("user not found", err)
			return status.Error(codes.NotFound, "User not found")
		}

		s.logger.Error("err DeleteUserByID ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Start making response")

	return nil

}
