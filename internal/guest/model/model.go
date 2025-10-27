package model

import "time"

type Guest struct {
	GuestID   int64      `gorm:"primaryKey;autoIncrement"`
	Name      string     `gorm:"type:varchar(500)"`
	Address   string     `gorm:"type:varchar(500)"`
	Phone     string     `gorm:"type:varchar(500)"`
	Email     string     `gorm:"type:varchar(500)"`
	EventId   string     `gorm:"type:varchar(500)"`
	CreatedAt *time.Time `gorm:"type:timestamp"`
	UpdatedAt *time.Time `gorm:"type:timestamp"`
	ProjectID int64      `gorm:"type:integer[]"`
}
