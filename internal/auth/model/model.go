package model

import "time"

// Auth maps to public.auth table
type Auth struct {
	UserID    int64      `gorm:"column:user_id"`
	Username  string     `gorm:"column:username;type:varchar(255)"`
	Password  string     `gorm:"column:password;type:text"`
	ProjectID int64      `gorm:"column:project_id;type:integer"`
	CreatedAt *time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:timestamp"`
}
