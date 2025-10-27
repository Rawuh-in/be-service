package db

import (
	"context"
	"errors"
	projectModel "rawuh-service/internal/project/model"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/model"
	"strconv"
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

func (p *ProjectRepository) ListProject(ctx context.Context, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort) (data []*projectModel.Project, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")

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

func (p *ProjectRepository) CreateProject(ctx context.Context, req *projectModel.CreateProjectRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")

	userInt, userIntErr := strconv.ParseInt(req.UserID, 0, 64)
	if userIntErr != nil {
		return userIntErr
	}

	now := time.Now()
	data := &projectModel.Project{
		ProjectName: req.ProjectName,
		CreatedById: userInt,
		CreatedAt:   &now,
		UpdatedById: userInt,
		UpdatedAt:   &now,
	}

	if err := query.Omit("event_id").Create(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) UpdateProject(ctx context.Context, req *projectModel.UpdateProjectRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")

	query = query.Where("project_id = ?", req.ProjectID)

	userInt, userIntErr := strconv.ParseInt(req.UserID, 0, 64)
	if userIntErr != nil {
		return userIntErr
	}

	now := time.Now()
	data := &projectModel.Project{
		ProjectName: req.ProjectName,
		UpdatedById: userInt,
		UpdatedAt:   &now,
	}

	if err := query.Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *ProjectRepository) DeleteProject(ctx context.Context, req *projectModel.DeleteProjectRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")

	query = query.Where("project_id = ?", req.ProjectID)

	if err := query.Delete(&projectModel.Project{}).Error; err != nil {
		return err
	}

	return nil

}

func (p *ProjectRepository) GetProjectDetail(ctx context.Context, req *projectModel.GetProjectDetailRequest) (*projectModel.Project, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.projects")

	query = query.Where("project_id = ?", req.ProjectID)

	var project projectModel.Project
	if err := query.First(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}
