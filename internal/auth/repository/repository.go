package repository

import (
	"context"
	"errors"

	model "rawuh-service/internal/auth/model"
	dbshared "rawuh-service/internal/shared/db"

	"gorm.io/gorm"
)

type AuthRepository struct {
	provider *dbshared.GormProvider
}

func NewAuthRepository(provider *dbshared.GormProvider) *AuthRepository {
	return &AuthRepository{provider: provider}
}

// GetAuthByUsername returns the auth row for the given username from public.auth
func (p *AuthRepository) GetAuthByUsername(ctx context.Context, username string) (*model.Auth, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var data model.Auth
	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.auth")
	query = query.Select("user_id, username, password, project_id, created_at, updated_at")
	query = query.Where("username = ?", username)

	if err := query.Take(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &data, nil
}

// CreateAuth inserts a new auth row into public.auth
func (p *AuthRepository) CreateAuth(ctx context.Context, a *model.Auth) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.auth")
	if err := query.Create(a).Error; err != nil {
		return err
	}
	return nil
}
