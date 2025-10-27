package db

import (
	"context"
	"errors"
	"time"

	guestModel "rawuh-service/internal/guest/model"
	"rawuh-service/internal/shared/db"
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

func (p *GuestRepository) CreateGuest(ctx context.Context, req *guestModel.CreateGuestRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	now := time.Now()
	data := &guestModel.Guest{
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		EventId:   req.EventId,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if err := query.Omit("guest_id").Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *GuestRepository) UpdateGuest(ctx context.Context, req *guestModel.UpdateGuestRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	query = query.Where("guest_id = ? and event_id = ?", req.GuestID, req.EventId)

	now := time.Now()
	data := &guestModel.Guest{
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
		Email:   req.Email,
		// EventId:   req.EventId,
		UpdatedAt: &now,
	}

	if err := query.Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *GuestRepository) GetGuestByID(ctx context.Context, guest_id string, event_id string) (*guestModel.Guest, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var data guestModel.Guest

	err := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.guests").
		Where("guest_id = ? and event_id = ?", guest_id, event_id).
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &guestModel.Guest{}, nil
		}

		return nil, err
	}

	return &data, nil
}

func (p *GuestRepository) DeleteGuestByID(ctx context.Context, guestID string, eventID string) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	db := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	tx := db.Where("guest_id = ? AND event_id = ?", guestID, eventID).Delete(&guestModel.Guest{})

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *GuestRepository) ListGuests(ctx context.Context, event_id string, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort) (data []*guestModel.Guest, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	query = query.Where("event_id = ?", event_id)

	query = query.Scopes(
		db.QueryScoop(sql.CollectiveAnd),
	)

	query = query.Scopes(db.Paginate(data, pagination, query))
	query = query.Scopes(
		db.Sort(sort),
	)

	if err := query.Debug().Find(&data).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return data, nil
}
