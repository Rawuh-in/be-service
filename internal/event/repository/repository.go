package db

import (
	"context"
	"errors"
	eventModel "rawuh-service/internal/event/model"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/model"
	"strconv"
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

func (p *EventRepository) ListEvent(ctx context.Context, participantID string, userID string, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort) (data []*eventModel.Event, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	query = query.Where("created_by_id = ? OR ? = ANY(participant_ids)", userID, participantID)

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

func (p *EventRepository) CreateEvent(ctx context.Context, req *eventModel.CreateEventRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	userInt, userIntErr := strconv.ParseInt(req.UserID, 0, 64)
	if userIntErr != nil {
		return userIntErr
	}

	now := time.Now()
	data := &eventModel.Event{
		EventName:   req.EventName,
		Description: req.Description,
		Options:     req.Options,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatedById: userInt,
		ProjectID:   userInt,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	if err := query.Omit("event_id").Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *EventRepository) UpdateEvent(ctx context.Context, req *eventModel.UpdateEventRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	query = query.Where("event_id = ?", req.EventID)

	userInt, userIntErr := strconv.ParseInt(req.UserID, 0, 64)
	if userIntErr != nil {
		return userIntErr
	}

	now := time.Now()
	data := &eventModel.Event{
		EventName:   req.EventName,
		Description: req.Description,
		Options:     req.Options,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		ProjectID:   userInt,
		UpdatedAt:   &now,
		UpdatedById: userInt,
	}

	if err := query.Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *EventRepository) GetEventByID(ctx context.Context, eventID string, userID string) (data *eventModel.Event, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	query = query.Where("event_id = ? AND (created_by_id = ? OR ? = ANY(participant_ids))", eventID, userID, userID)

	if err := query.Debug().Find(&data).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	return data, nil

}

func (p *EventRepository) DeleteEventByID(ctx context.Context, eventID string, userID string) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.events")

	tx := query.Where("event_id = ? AND (created_by_id = ? OR ? = ANY(participant_ids))", eventID, userID, userID)

	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

}
