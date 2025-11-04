package model

import "time"

type User struct {
	UserID        int64      `gorm:"primaryKey;autoIncrement"`
	Name          string     `gorm:"type:varchar(500)"`
	UserType      string     `gorm:"type:varchar(100)"`
	Username      string     `gorm:"type:varchar(500)"`
	Email         string     `gorm:"type:varchar(500)"`
	ProjectID     int64      `gorm:"type:integer"`
	CreatedById   int64      `gorm:"type:integer"`
	CreatedByName string     `gorm:"type:varchar(500)"`
	UpdatedById   int64      `gorm:"type:integer"`
	UpdatedByName string     `gorm:"type:varchar(500)"`
	EventId       int64      `gorm:"type:integer"`
	CreatedAt     *time.Time `gorm:"type:timestamp"`
	UpdatedAt     *time.Time `gorm:"type:timestamp"`
	Status        int64      `gorm:"type:integer"`
}
