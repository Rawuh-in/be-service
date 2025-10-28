package config

import (
	"rawuh-service/internal/shared/lib/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Dsn string `config:"DB_DSN"`

	// RedisAddr string `config:"REDIS_ADDR"`
	// RedisPass string `config:"REDIS_PASS"`
	// RedisDB   int    `config:"REDIS_DB"`
}

var Configs *Config

func InitConfig() *Config {
	return &Config{
		Dsn: utils.GetEnv("DB_DSN", ""),
		// RedisAddr: utils.GetEnv("REDIS_ADDR", ""),
		// RedisPass: utils.GetEnv("REDIS_PASS", ""),
		// RedisDB:   0,
	}
}

func InitDB(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Optional: set connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil // ✅ now returning *gorm.DB
}

// func NewLogger() *logrus.Logger {
// 	log := logrus.New()
// 	log.SetFormatter(&logrus.JSONFormatter{})
// 	log.SetLevel(logrus.InfoLevel)
// 	return log
// }
