package db

import (
	"context"
	"errors"
	"time"

	guestModel "rawuh-service/internal/guest/model"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/middleware"
	model "rawuh-service/internal/shared/model"

	"gorm.io/gorm"
)

type GuestRepository struct {
	provider *db.GormProvider
}

func NewGuestRepository(provider *db.GormProvider) *GuestRepository {
	return &GuestRepository{
		provider: provider,
	}
}

func (p *GuestRepository) CreateGuest(ctx context.Context, req *guestModel.CreateGuestRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	now := time.Now()
	data := &guestModel.Guest{
		ProjectID: currentUser.ProjectID,
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		EventId:   currentUser.EventID,
		CreatedAt: &now,
		Options:   req.Options,
	}

	if err := query.Omit("guest_id").Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *GuestRepository) UpdateGuest(ctx context.Context, req *guestModel.UpdateGuestRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	query = query.Where("project_id = ? and guest_id = ? and event_id = ?", currentUser.ProjectID, req.GuestID, currentUser.EventID)

	now := time.Now()
	data := &guestModel.Guest{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
		Options: req.Options,
		// EventId:   req.EventId,
		UpdatedAt: &now,
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

func (p *GuestRepository) GetGuestByID(ctx context.Context, guestID string, currentUser middleware.AuthClaims) (*guestModel.Guest, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var data guestModel.Guest

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.guests")

	query = query.Where("project_id = ? and guest_id = ? and event_id = ?", currentUser.ProjectID, guestID, currentUser.EventID)

	if err := query.Debug().Find(&data).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

	}

	return &data, nil
}

func (p *GuestRepository) DeleteGuestByID(ctx context.Context, guestID string, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	query = query.Where("project_id = ? AND guest_id = ? AND event_id = ?", currentUser.ProjectID, guestID, currentUser.EventID)

	res := query.Delete(&guestModel.Guest{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *GuestRepository) ListGuests(ctx context.Context, currentUser middleware.AuthClaims, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort) (data []*guestModel.Guest, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	query = query.Where("project_id = ? AND event_id = ?", currentUser.ProjectID, currentUser.EventID)

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
