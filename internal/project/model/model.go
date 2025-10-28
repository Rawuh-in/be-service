package model

import "time"

type Project struct {
	ProjectID   int64      `gorm:"primaryKey;autoIncrement"`
	ProjectName string     `gorm:"type:varchar(500)"`
	CreatedAt   *time.Time `gorm:"type:timestamp"`
	CreatedById int64      `gorm:"type:bigint"`
	UpdatedAt   *time.Time `gorm:"type:timestamp"`
	UpdatedById int64      `gorm:"type:bigint"`
	Status      string     `gorm:"type:varchar(100)"`
	StatusDesc  string     `gorm:"type:varchar(500)"`
}
