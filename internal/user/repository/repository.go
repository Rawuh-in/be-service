package db

import (
	"context"
	"errors"
	"strconv"
	"time"

	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/middleware"
	model "rawuh-service/internal/shared/model"
	userModel "rawuh-service/internal/user/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	provider *db.GormProvider
}

func NewUserRepository(provider *db.GormProvider) *UserRepository {
	return &UserRepository{
		provider: provider,
	}
}

func (p *UserRepository) CreateUser(ctx context.Context, req *userModel.CreateUserRequest, currentUser middleware.AuthClaims) (int64, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.users")

	projectID, _ := strconv.ParseInt(req.ProjectID, 0, 64)
	// userID, _ := strconv.ParseInt(req.UserID, 0, 64)

	now := time.Now()
	data := &userModel.User{
		ProjectID:     projectID,
		Name:          req.Name,
		UserType:      req.UserType,
		Username:      req.Username,
		Email:         req.Email,
		CreatedById:   currentUser.UserID,
		CreatedByName: currentUser.Name,
		Status:        1,
		CreatedAt:     &now,
	}

	if err := query.Omit("user_id").Create(data).Error; err != nil {
		return 0, err
	}

	return data.UserID, nil
}

func (p *UserRepository) UpdateUser(ctx context.Context, req *userModel.UpdateUserRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.users")

	query = query.Where("user_id = ?", req.UserID)

	projectID, _ := strconv.ParseInt(req.ProjectID, 0, 64)
	// eventID, _ := strconv.ParseInt(req.EventId, 0, 64)
	// userID, _ := strconv.ParseInt(req.UserID, 0, 64)

	now := time.Now()
	data := &userModel.User{
		ProjectID:     projectID,
		Name:          req.Name,
		UserType:      req.UserType,
		Username:      req.Username,
		Email:         req.Email,
		UpdatedById:   currentUser.UserID,
		UpdatedByName: currentUser.Name,
		UpdatedAt:     &now,
	}

	res := query.Updates(data)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *UserRepository) GetUserByID(ctx context.Context, userID string) (*userModel.User, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var data userModel.User

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.users")

	query = query.Where("user_id = ?", userID)

	if err := query.Debug().Find(&data).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return &data, nil
}

func (p *UserRepository) DeleteUserByID(ctx context.Context, userID string) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.users")

	query = query.Where("user_id = ?", userID)

	res := query.Delete(&userModel.User{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *UserRepository) ListUsers(ctx context.Context, projectID string, eventID string, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort) (data []*userModel.User, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.users")

	query = query.Scopes(
		db.QueryScoop(sql.CollectiveAnd),
	)

	query = query.Scopes(db.Paginate(data, pagination, query))
	query = query.Scopes(
		db.Sort(sort),
	)

	if err := query.Debug().First(&data).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return data, nil
}

func (p *UserRepository) CheckUsernameExist(ctx context.Context, username string) (bool, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var data userModel.User

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.users")

	query = query.Where("username = ?", username)

	if err := query.Debug().Find(&data).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
	}

	return data.UserID != 0, nil
}
