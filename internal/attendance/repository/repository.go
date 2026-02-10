package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	attendanceModel "rawuh-service/internal/attendance/model"
	guestModel "rawuh-service/internal/guest/model"
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

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.attendances")
	query = query.Where("project_id = ? AND event_id = ?", req.ProjectID, req.EventId)

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

func (p *AttendanceRepository) GetAttendanceByID(ctx context.Context, req *attendanceModel.GetAttendanceByIDRequest) (*attendanceModel.Attendance, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var data attendanceModel.Attendance

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.attendances")

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

	fmt.Println("START CREATE ATTENDANCE")

	eventInt, _ := strconv.ParseInt(req.EventId, 0, 64)
	projectInt, _ := strconv.ParseInt(req.ProjectID, 0, 64)

	now := time.Now()
	var attendances []*attendanceModel.Attendance

	for _, guest := range req.GuestData {
		// guestInt, _ := strconv.ParseInt(*guest.GuestID, 0, 64)
		attendances = append(attendances, &attendanceModel.Attendance{
			ProjectID:     projectInt,
			EventID:       eventInt,
			GuestID:       guest.GuestID,
			GuestName:     guest.Name,
			Status:        1,
			StatusStr:     "Active",
			CreatedAt:     &now,
			CreatedByName: currentUser.Name,
			CreatedByID:   fmt.Sprint(currentUser.UserID),
		})
	}

	if err := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.attendances").
		Omit("ID").
		CreateInBatches(attendances, 100).Error; err != nil {
		return err
	}

	return nil
}

func (p *AttendanceRepository) UpdateAttendance(ctx context.Context, req *attendanceModel.UpdateAttendanceRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.attendances")
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

func (p *AttendanceRepository) CheckInOutAttendance(ctx context.Context, req *attendanceModel.BulkCheckInOutAttendanceRequest, currentUser middleware.AuthClaims) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	now := time.Now()

	updates := &attendanceModel.Attendance{
		UpdatedAt: &now,
	}

	selectFields := []string{"updated_at", "status", "status_str"}

	switch req.Type {
	case constant.AttendanceTypeCheckIn:
		updates.CheckedInAt = &now
		updates.Status = 2
		updates.StatusStr = "Checked In"
		selectFields = append(selectFields, "checked_in_at")
	case constant.AttendanceTypeCheckOut:
		updates.CheckedOutAt = &now
		updates.Status = 3
		updates.StatusStr = "Checked Out"
		selectFields = append(selectFields, "checked_out_at")
	}

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.attendances").
		Select(selectFields).
		Where("project_id = ? AND event_id = ? AND attendance_id IN ?", req.ProjectID, req.EventID, req.AttendanceIDs)

	if err := query.Updates(updates).Error; err != nil {
		return err
	}

	return nil
}

func (p *AttendanceRepository) GetListGuestByEventID(ctx context.Context, currentUser middleware.AuthClaims, req *attendanceModel.CreateAttendanceRequest) ([]*guestModel.Guest, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.guests")

	query = query.Where("project_id = ? AND event_id = ?", req.ProjectID, req.EventId)

	var guests []*guestModel.Guest
	if err := query.Find(&guests).Error; err != nil {
		return nil, err
	}

	return guests, nil
}

func (p *AttendanceRepository) GetAttendancesByIDs(ctx context.Context, req *attendanceModel.DeleteAttendanceByIDRequest) ([]*attendanceModel.Attendance, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var attendances []*attendanceModel.Attendance

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.attendances").
		Where("attendance_id IN ? AND project_id = ? AND event_id = ?", req.AttendanceIDs, req.ProjectID, req.EventID)

	if err := query.Find(&attendances).Error; err != nil {
		return nil, err
	}

	return attendances, nil
}

func (p *AttendanceRepository) DeleteAttendanceByIDs(ctx context.Context, req *attendanceModel.DeleteAttendanceByIDRequest) error {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.attendances")

	query = query.Where("project_id = ? AND event_id = ? AND attendance_id IN ?", req.ProjectID, req.EventID, req.AttendanceIDs)

	res := query.Delete(&attendanceModel.Attendance{})

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *AttendanceRepository) GetListAttendanceByID(ctx context.Context, req *attendanceModel.BulkCheckInOutAttendanceRequest) ([]string, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.attendances")

	query = query.Where("project_id = ? AND event_id = ? and attendance_id in (?)", req.ProjectID, req.EventID, req.AttendanceIDs)
	var attendanceIDsResult []string
	if err := query.Select("attendance_id").Find(&attendanceIDsResult).Error; err != nil {
		return nil, err
	}

	return attendanceIDsResult, nil
}

func (p *AttendanceRepository) GetExistingAttendanceByGuestIDs(ctx context.Context, guestIDs []int64, eventID int64, projectID int64) ([]int64, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().Table("public.attendances")

	query = query.Where("project_id = ? AND event_id = ? AND guest_id IN ?", projectID, eventID, guestIDs)

	var existingGuestIDs []int64
	if err := query.Select("DISTINCT guest_id").Find(&existingGuestIDs).Error; err != nil {
		return nil, err
	}

	return existingGuestIDs, nil
}

func (p *AttendanceRepository) GetAttendanceRecordsByIDs(ctx context.Context, attendanceIDs []string, projectID, eventID string) ([]*attendanceModel.Attendance, error) {
	timeoutctx, cancel := context.WithTimeout(ctx, p.provider.GetTimeout())
	defer cancel()

	var attendances []*attendanceModel.Attendance

	query := p.provider.GetDB().WithContext(timeoutctx).Debug().
		Table("public.attendances").
		Where("attendance_id IN ? AND project_id = ? AND event_id = ?", attendanceIDs, projectID, eventID)

	if err := query.Find(&attendances).Error; err != nil {
		return nil, err
	}

	return attendances, nil
}
