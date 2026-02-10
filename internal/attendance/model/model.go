package model

import "time"

type Attendance struct {
	ID            int64      `gorm:"primaryKey;autoIncrement;column:attendance_id"`
	EventID       int64      `gorm:"index;not null;uniqueIndex:idx_event_guest"`
	GuestID       int64      `gorm:"index;not null;uniqueIndex:idx_event_guest"`
	ProjectID     int64      `gorm:"index;not null;uniqueIndex:idx_event_guest"`
	GuestName     string     `gorm:"type:varchar(255)"`
	CheckedInAt   *time.Time `gorm:"type:timestamp"`
	CheckedOutAt  *time.Time `gorm:"type:timestamp"`
	Status        int64      `gorm:"type:varchar(50)"`
	StatusStr     string     `gorm:"type:varchar(50)"`
	CreatedAt     *time.Time `gorm:"type:timestamp"`
	CreatedByName string     `gorm:"type:varchar(255)"`
	CreatedByID   string     `gorm:"type:varchar(255);index"`
	UpdatedByName string     `gorm:"type:varchar(255)"`
	UpdatedByID   string     `gorm:"type:varchar(255);index"`
	UpdatedAt     *time.Time `gorm:"type:timestamp"`
}
