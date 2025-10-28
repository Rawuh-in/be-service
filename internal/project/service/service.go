package service

import (
	"context"
	"encoding/base64"
	"net/http"
	projectModel "rawuh-service/internal/project/model"
	projectDb "rawuh-service/internal/project/repository"
	"rawuh-service/internal/shared/constant"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/logger"
	"rawuh-service/internal/shared/model"
	"rawuh-service/internal/shared/redis"
	"strconv"
	"strings"

	"go.elastic.co/apm/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectService interface {
	ListProjects(ctx context.Context, req *projectModel.ListProjectRequest) (*projectModel.ListProjectResponse, error)
	CreateProject(ctx context.Context, req *projectModel.CreateProjectRequest) error
	UpdateProject(ctx context.Context, req *projectModel.UpdateProjectRequest) error
	DeleteProject(ctx context.Context, req *projectModel.DeleteProjectRequest) error
	GetProjectDetail(ctx context.Context, req *projectModel.GetProjectDetailRequest) (*projectModel.GetProjectDetailResponse, error)
}

type projectService struct {
	dbProvider *projectDb.ProjectRepository
	logger     *logger.Logger
	redis      *redis.Redis
}

func NewProjectService(dbProvider *projectDb.ProjectRepository, logger *logger.Logger, redis *redis.Redis) ProjectService {
	return &projectService{
		dbProvider: dbProvider,
		logger:     logger,
		redis:      redis,
	}
}

func (s *projectService) ListProjects(ctx context.Context, req *projectModel.ListProjectRequest) (*projectModel.ListProjectResponse, error) {
	funcName := "ListProjects"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)

	// if req.EventId == "" {
	// 	s.logger.Error("err Invalid event id : ", req.EventId)
	// 	return nil, status.Errorf(codes.InvalidArgument, "Invalid Event Id")
	// }

	decodeQuery, err := base64.RawStdEncoding.DecodeString(req.Query)
	if err != nil {
		s.logger.Error("err DecodeString ", err)
		return nil, nil
	}

	loggerZap.Info("Success Decode Query")

	pagination := utils.SetPagination(req.Page, req.Limit)

	allowedColumns := map[string]bool{
		"created_at":   true,
		"updated_at":   true,
		"status":       true,
		"project_name": true,
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

	loggerZap.Info("Start ListProjects with data ", req)
	projects, err := s.dbProvider.ListProject(ctx, pagination, sqlBuilder, sort)
	if err != nil {
		loggerZap.Error("err ListProjects ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	loggerZap.Info("Start making response")

	result := &projectModel.ListProjectResponse{
		Error:      false,
		Code:       http.StatusOK,
		Message:    "Success",
		Data:       projects,
		Pagination: pagination,
	}

	return result, nil
}

func (s *projectService) CreateProject(ctx context.Context, req *projectModel.CreateProjectRequest) error {
	funcName := "CreateProject"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	loggerZap.Debug("Start GetMeFromMD")

	nameLength, _ := strconv.Atoi(utils.GetEnv("PRODUCT_NAME_LENGTH", "255"))

	s.logger.Info("Start CreateProject Validation for req ", req)

	if utils.IsEmptyString(req.ProjectName) {
		return status.Errorf(codes.Aborted, "project name is empty")
	}
	if len(req.ProjectName) > nameLength {
		return status.Errorf(codes.Aborted, "project name maximum characters is %d", nameLength)
	}

	if !utils.IsValidProductName(req.ProjectName) {
		return status.Errorf(codes.Aborted, "characters not allowed in project name")
	}

	if req.UserID == "" {
		return status.Errorf(codes.Aborted, "user id is empty")
	}

	loggerZap.Info("Start CreateProject with data ", req)

	err := s.dbProvider.CreateProject(ctx, req)
	if err != nil {
		s.logger.Error("err CreateGuest ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Success CreateGuest")

	return nil
}

func (s *projectService) UpdateProject(ctx context.Context, req *projectModel.UpdateProjectRequest) error {
	funcName := "UpdateProject"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)
	loggerZap.Debug("Start GetMeFromMD")

	nameLength, _ := strconv.Atoi(utils.GetEnv("PRODUCT_NAME_LENGTH", "255"))

	s.logger.Info("Start UpdateProject Validation for req ", req)

	if utils.IsEmptyString(req.ProjectName) {
		return status.Errorf(codes.Aborted, "project name is empty")
	}
	if len(req.ProjectName) > nameLength {
		return status.Errorf(codes.Aborted, "project name maximum characters is %d", nameLength)
	}

	// if req.Status == 0 {
	// 	return status.Errorf(codes.Aborted, "status is empty")
	// }

	// if utils.IsEmptyString(req.StatusDesc) {
	// 	return status.Errorf(codes.Aborted, "status description is empty")
	// }

	loggerZap.Info("Start UpdateProject with data ", req)

	err := s.dbProvider.UpdateProject(ctx, req)
	if err != nil {
		s.logger.Error("err UpdateProject ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Success UpdateProject")

	return nil
}

func (s *projectService) DeleteProject(ctx context.Context, req *projectModel.DeleteProjectRequest) error {
	funcName := "DeleteProject"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)

	if req.ProjectID == "" {
		loggerZap.Error("invalid project id", nil)
		return status.Errorf(codes.Aborted, "project id is empty")
	}

	loggerZap.Info("Start DeleteProject with data ", req)

	err := s.dbProvider.DeleteProject(ctx, req)
	if err != nil {
		s.logger.Error("err DeleteProject ", err)
		return status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Success DeleteProject")

	return nil
}

func (s *projectService) GetProjectDetail(ctx context.Context, req *projectModel.GetProjectDetailRequest) (*projectModel.GetProjectDetailResponse, error) {
	funcName := "GetProjectDetail"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, req)

	if req.ProjectID == "" {
		loggerZap.Error("invalid project id", nil)
		return nil, status.Errorf(codes.Aborted, "project id is empty")
	}

	loggerZap.Info("Start GetProjectDetail with data ", req)

	project, err := s.dbProvider.GetProjectDetail(ctx, req)
	if err != nil {
		s.logger.Error("err GetProjectDetail ", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	s.logger.Info("Success GetProjectDetail")

	result := &projectModel.GetProjectDetailResponse{
		Error:   false,
		Code:    http.StatusOK,
		Message: "Success",
		Data:    project,
	}

	return result, nil
}
