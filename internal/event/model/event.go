package event_model

import (
	"time"
)

type Events struct {
	EventID       int64      `gorm:"primaryKey;autoIncrement"`
	EventName     string     `gorm:"type:varchar(500)"`
	Description   string     `gorm:"type:varchar(500)"`
	StartDate     *time.Time `gorm:"type:timestamp"`
	EndDate       *time.Time `gorm:"type:timestamp"`
	Options       string     `gorm:"type:jsonb"`
	CreatedAt     *time.Time `gorm:"type:timestamp"`
	UpdatedAt     *time.Time `gorm:"type:timestamp"`
	ProjectID     int64      `gorm:"type:integer[]"`
	CreatedById   int64      `gorm:"type:bigint"`
	CreatedByName string     `gorm:"type:varchar(500)"`
}
