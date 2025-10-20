package db

import (
	"fmt"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/model"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// var location, _ = time.LoadLocation("Asia/Jakarta")

type QueryBuilder struct {
	CollectiveAnd string
	CustomOrder   string
	Sort          *model.Sort
}

type GormProvider struct {
	db_main *gorm.DB
	timeout time.Duration
}

func NewProvider(db *gorm.DB) *GormProvider {
	timeoutDB := utils.GetEnv("DB_TIMEOUT", "300")
	timeoutduration, err := strconv.Atoi(timeoutDB)
	if err != nil {
		fmt.Println("", err)
	}

	return &GormProvider{db_main: db, timeout: time.Duration(timeoutduration) * time.Second}
}

func (p *GormProvider) NewTransaction() *gorm.DB {
	return p.db_main.Begin()
}

// GetDB returns the underlying *gorm.DB
func (p *GormProvider) GetDB() *gorm.DB {
	return p.db_main
}

// GetTimeout returns the configured timeout
func (p *GormProvider) GetTimeout() time.Duration {
	return p.timeout
}
