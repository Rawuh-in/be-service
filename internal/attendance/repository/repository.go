package db

import (
	"context"
	"errors"
	"strconv"
	"time"

	attendanceModel "rawuh-service/internal/attendance/model"
	"rawuh-service/internal/shared/constant"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/middleware"
	model "rawuh-service/internal/shared/model"

	"gorm.io/gorm"
)

type AttendanceRepository struct {
	provider *db.GormProvider
}

func NewAttendanceRepository(provider *db.GormProvider) *AttendanceRepository {
	return &AttendanceRepository{
		provider: provider,
	}
}

func (p *AttendanceRepository) ListAttendance(ctx context.Context, req *attendanceModel.ListAttendanceRequest, pagination *model.PaginationResponse, sql *db.QueryBuilder, sort *model.Sort) (data []*attendanceModel.Attendance, err error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.attendance")
	query = query.Where("project_id = ? AND event_id = ?", req.ProjectID, req.EventId)

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

func (p *AttendanceRepository) GetAttendanceByID(ctx context.Context, req *attendanceModel.GetAttendanceByIDRequest) (*attendanceModel.Attendance, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var data attendanceModel.Attendance

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.attendance")

	query = query.Where("project_id = ? and attendance_id = ? and event_id = ?", req.ProjectID, req.AttendanceID, req.EventId)

	if err := query.Debug().Find(&data).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}

	}

	return &data, nil
}

func (p *AttendanceRepository) CreateAttendance(ctx context.Context, req *attendanceModel.CreateAttendanceRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	eventInt, _ := strconv.ParseInt(req.EventId, 0, 64)
	projectInt, _ := strconv.ParseInt(req.ProjectID, 0, 64)

	now := time.Now()
	var attendances []*attendanceModel.Attendance

	for _, guestID := range req.GuestID {
		guestInt, _ := strconv.ParseInt(*guestID, 0, 64)
		attendances = append(attendances, &attendanceModel.Attendance{
			ProjectID: projectInt,
			EventID:   eventInt,
			GuestID:   guestInt,
			CreatedAt: &now,
		})
	}

	if err := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.attendance").
		Omit("id").
		CreateInBatches(attendances, 100).Error; err != nil {
		return err
	}

	return nil
}

func (p *AttendanceRepository) UpdateAttendance(ctx context.Context, req *attendanceModel.UpdateAttendanceRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.attendance")
	query = query.Where("project_id = ? and attendance_id = ? and event_id = ?", req.ProjectID, req.AttendanceID, req.EventId)

	eventInt, _ := strconv.ParseInt(req.EventId, 0, 64)
	projectInt, _ := strconv.ParseInt(req.ProjectID, 0, 64)
	guestInt, _ := strconv.ParseInt(req.GuestID, 0, 64)

	now := time.Now()
	data := &attendanceModel.Attendance{
		ProjectID: projectInt,
		EventID:   eventInt,
		GuestID:   guestInt,
		UpdatedAt: &now,
	}

	switch req.Type {
	case constant.AttendanceTypeCheckIn:
		data.CheckedInAt = &now
	case constant.AttendanceTypeCheckOut:
		data.CheckedOutAt = &now
	}

	if err := query.Updates(data).Error; err != nil {
		return err
	}

	return nil
}

func (p *AttendanceRepository) CheckInOutAttendance(ctx context.Context, req *attendanceModel.CheckInOutAttendanceRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	now := time.Now()

	updates := &attendanceModel.Attendance{
		UpdatedAt: &now,
	}

	switch req.Type {
	case constant.AttendanceTypeCheckIn:
		updates.CheckedInAt = &now
	case constant.AttendanceTypeCheckOut:
		updates.CheckedOutAt = &now
	}

	// Bulk update all attendance records matching the attendance IDs
	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.attendance").
		Where("project_id = ? AND event_id = ? AND attendance_id IN ?", req.ProjectID, req.EventID, req.AttendanceIDs)

	if err := query.Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (p *AttendanceRepository) GetListGuestByEventID(ctx context.Context, currentUser middleware.AuthClaims) ([]string, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	eventID := currentUser.EventID
	projectID := currentUser.ProjectID

	query = query.Where("project_id = ? AND event_id = ?", projectID, eventID)

	var guestIDs []string
	if err := query.Select("guest_id").Find(&guestIDs).Error; err != nil {
		return nil, err
	}

	return guestIDs, nil
}

func (p *AttendanceRepository) GetAttendancesByIDs(ctx context.Context, attendanceIDs []string, currentUser middleware.AuthClaims) ([]*attendanceModel.Attendance, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var attendances []*attendanceModel.Attendance

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.attendance").
		Where("attendance_id IN ? AND project_id = ? AND event_id = ?", attendanceIDs, currentUser.ProjectID, currentUser.EventID)

	if err := query.Find(&attendances).Error; err != nil {
		return nil, err
	}

	return attendances, nil
}

func (p *AttendanceRepository) DeleteAttendanceByIDs(ctx context.Context, attendanceIDs []string, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.attendance")

	query = query.Where("project_id = ? AND event_id = ? AND attendance_id IN ?", currentUser.ProjectID, currentUser.EventID, attendanceIDs)

	res := query.Delete(&attendanceModel.Attendance{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
