package service

import (
	"context"
	"rawuh-service/internal/auth/model"
	"rawuh-service/internal/auth/repository"
	"rawuh-service/internal/shared/constant"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/logger"

	"go.elastic.co/apm/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService interface {
	Authenticate(ctx context.Context, username, password string) (*model.Auth, error)
}

type authService struct {
	repo   *repository.AuthRepository
	logger *logger.Logger
}

func NewAuthService(repo *repository.AuthRepository, logger *logger.Logger) AuthService {
	return &authService{repo: repo, logger: logger}
}

func (s *authService) Authenticate(ctx context.Context, username, password string) (*model.Auth, error) {
	funcName := "AuthService.Authenticate"
	span, ctx := apm.StartSpan(ctx, funcName, constant.SpanTypeProccess)
	span.Action = constant.SpanActionExecute
	defer span.End()

	ctx, loggerZap := s.logger.StartLogger(ctx, funcName, map[string]interface{}{"username": username})

	if username == "" || password == "" {
		loggerZap.Warn("empty username or password")
		return nil, status.Error(codes.InvalidArgument, "username or password empty")
	}

	auth, err := s.repo.GetAuthByUsername(ctx, username)
	if err != nil {
		loggerZap.Error("err GetAuthByUsername", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}
	if auth == nil {
		loggerZap.Info("auth not found for username")
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	decrypted, err := utils.DecryptAES(auth.Password)
	if err != nil {
		loggerZap.Error("err DecryptAES", err)
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	if decrypted != password {
		loggerZap.Warn("invalid password for user", nil)
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	loggerZap.Info("authentication success for user")
	return auth, nil
}
