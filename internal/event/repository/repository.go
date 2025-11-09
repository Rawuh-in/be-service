package db

import (
	"context"
	"errors"
	eventModel "rawuh-service/internal/event/model"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/middleware"
	"rawuh-service/internal/shared/model"
	"time"

	"gorm.io/gorm"
)

type EventRepository struct {
	provider *db.GormProvider
}

func NewEventRepository(provider *db.GormProvider) *EventRepository {
	return &EventRepository{
		provider: provider,
	}
}

func (p *EventRepository) ListEvent(ctx context.Context, currentUser middleware.AuthClaims, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort) (data []*eventModel.Event, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	query = query.Where("project_id = ?", currentUser.ProjectID)

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

func (p *EventRepository) CreateEvent(ctx context.Context, req *eventModel.CreateEventRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	now := time.Now()
	data := &eventModel.Event{
		EventName:    req.EventName,
		Description:  req.Description,
		EventOptions: req.EventOptions,
		GuestOptions: req.GuestOptions,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		CreatedById:  currentUser.UserID,
		ProjectID:    currentUser.ProjectID,
		CreatedAt:    &now,
	}

	if err := query.Omit("event_id").Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *EventRepository) UpdateEvent(ctx context.Context, req *eventModel.UpdateEventRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	query = query.Where("event_id = ?", req.EventID)

	now := time.Now()
	data := &eventModel.Event{
		EventName:    req.EventName,
		Description:  req.Description,
		EventOptions: req.EventOptions,
		GuestOptions: req.GuestOptions,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		ProjectID:    currentUser.ProjectID,
		UpdatedAt:    &now,
		UpdatedById:  currentUser.UserID,
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

func (p *EventRepository) GetEventByID(ctx context.Context, eventID string, currentUser middleware.AuthClaims) (data *eventModel.Event, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	query = query.Where("project_id = ? and event_id = ?", currentUser.ProjectID, eventID)

	if err := query.Debug().First(&data).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return data, nil

}

func (p *EventRepository) DeleteEventByID(ctx context.Context, eventID string, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	query = query.Where("project_id = ? and event_id = ?", currentUser.ProjectID, eventID)

	res := query.Delete(&eventModel.Event{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}
