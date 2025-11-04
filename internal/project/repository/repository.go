package db

import (
	"context"
	"errors"
	projectModel "rawuh-service/internal/project/model"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/middleware"
	"rawuh-service/internal/shared/model"
	"time"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	provider *db.GormProvider
}

func NewProjectRepository(provider *db.GormProvider) *ProjectRepository {
	return &ProjectRepository{
		provider: provider,
	}
}

func (p *ProjectRepository) ListProject(ctx context.Context, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort, projectID int64) (data []*projectModel.Project, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")

	query = query.Where("project_id = ?", projectID)

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

func (p *ProjectRepository) CreateProject(ctx context.Context, req *projectModel.CreateProjectRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")

	now := time.Now()
	data := &projectModel.Project{
		ProjectName: req.ProjectName,
		CreatedById: currentUser.UserID,
		CreatedAt:   &now,
		Status:      1,
	}

	if err := query.Omit("event_id").Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) UpdateProject(ctx context.Context, req *projectModel.UpdateProjectRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")

	query = query.Where("project_id = ?", req.ProjectID)

	now := time.Now()
	data := &projectModel.Project{
		ProjectName: req.ProjectName,
		UpdatedById: currentUser.UserID,
		UpdatedAt:   &now,
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

func (p *ProjectRepository) DeleteProject(ctx context.Context, req *projectModel.DeleteProjectRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")
	query = query.Where("project_id = ?", req.ProjectID)

	res := query.Delete(&projectModel.Project{})
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *ProjectRepository) GetProjectDetail(ctx context.Context, req *projectModel.GetProjectDetailRequest) (*projectModel.Project, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")

	query = query.Where("project_id = ?", req.ProjectID)

	var project projectModel.Project
	if err := query.Debug().Find(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &project, nil
}
