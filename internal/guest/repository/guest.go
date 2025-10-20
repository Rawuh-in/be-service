package guest_db

import (
	"context"
	"errors"
	"time"

	guest_model "rawuh-service/internal/guest/model"
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

// func (p *GuestRepository) InsertProduct(ctx context.Context, req *guest_model.CreateProductRequest) error {
// 	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
// 	defer cancel()

// 	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.products")

// 	now := time.Now()
// 	data := &guest_model.Product{
// 		Name:        req.Name,
// 		Price:       req.Price,
// 		Description: req.Description,
// 		Quantity:    req.Quantity,
// 		CreatedAt:   &now,
// 		UpdatedAt:   &now,
// 	}

// 	if err := query.Create(data).Error; err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (p *GuestRepository) GetProductByName(ctx context.Context, productName string) (*guest_model.Product, error) {
// 	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
// 	defer cancel()

// 	var data guest_model.Product

// 	err := p.provider.GetDB().WithContext(timeoutctx).
// 		Table("public.products").
// 		Where("name = ?", productName).
// 		First(&data).Error

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return &guest_model.Product{}, nil
// 		}

// 		return nil, err
// 	}

// 	return &data, nil
// }

// func (p *GuestRepository) ListProduct(ctx context.Context, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort) (data []*guest_model.Product, err error) {
// 	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
// 	defer cancel()

// 	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.products")

// 	query = query.Scopes(
// 		db.QueryScoop(sql.CollectiveAnd),
// 	)

// 	query = query.Scopes(db.Paginate(data, pagination, query))
// 	query = query.Scopes(
// 		db.Sort(sort),
// 	)

// 	if err := query.Debug().Find(&data).Error; err != nil {
// 		if !errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, err
// 		}
// 	}

// 	return data, nil
// }

func (p *GuestRepository) CreateGuest(ctx context.Context, req *guest_model.CreateGuestRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	now := time.Now()
	data := &guest_model.Guests{
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Email:     req.Email,
		EventId:   req.EventId,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if err := query.Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *GuestRepository) GetGuestByName(ctx context.Context, guestName string) (*guest_model.Guests, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var data guest_model.Guests

	err := p.provider.GetDB().WithContext(timeoutctx).
		Table("public.guests").
		Where("name = ?", guestName).
		First(&data).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &guest_model.Guests{}, nil
		}

		return nil, err
	}

	return &data, nil
}

func (p *GuestRepository) ListGuests(ctx context.Context, event_id string, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort) (data []*guest_model.Guests, err error) {
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
