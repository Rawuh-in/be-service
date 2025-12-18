package model

import "time"

type Attendance struct {
	ID           int64      `gorm:"primaryKey;autoIncrement"`
	EventID      int64      `gorm:"index;not null;uniqueIndex:idx_event_guest"`
	GuestID      int64      `gorm:"index;not null;uniqueIndex:idx_event_guest"`
	ProjectID    int64      `gorm:"index;not null;uniqueIndex:idx_event_guest"`
	CheckedInAt  *time.Time `gorm:"type:timestamp"`
	CheckedOutAt *time.Time `gorm:"type:timestamp"`
	Status       int64      `gorm:"type:varchar(50);default:'not_checked_in'"`
	StatusStr    string     `gorm:"type:varchar(50);default:'not_checked_in'"`
	CreatedAt    *time.Time `gorm:"type:timestamp"`
	UpdatedAt    *time.Time `gorm:"type:timestamp"`
}
