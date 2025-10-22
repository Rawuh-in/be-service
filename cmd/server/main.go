package main

import (
	"log"
	"net/http"
	eventHandler "rawuh-service/internal/event/handler"
	eventDb "rawuh-service/internal/event/repository"
	eventService "rawuh-service/internal/event/service"
	guestHandler "rawuh-service/internal/guest/handler"
	guestDb "rawuh-service/internal/guest/repository"
	guestService "rawuh-service/internal/guest/service"
	"rawuh-service/internal/shared/config"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/redis"
	"rawuh-service/internal/shared/router"

	"github.com/joho/godotenv"
)

var appConfig *config.Config

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	appConfig = config.InitConfig()
	logger := config.NewLogger()

	logger.Info("Start connecting to db ", appConfig.Dsn)

	gormDB, err := config.InitDB(appConfig.Dsn)
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	logger.Info("Success connecting to db ", appConfig.Dsn)

	dbProvider := db.NewProvider(gormDB)
	guestDB := guestDb.NewGuestRepository(dbProvider)
	eventDB := eventDb.NewEventRepository(dbProvider)

	rdb := redis.NewRedis(appConfig.RedisAddr, appConfig.RedisPass, appConfig.RedisDB)

	guestService := guestService.NewGuestService(guestDB, logger, rdb)
	guestHandler := guestHandler.NewGuestHandler(guestService)

	eventService := eventService.NewEventService(eventDB, logger, rdb)
	eventHandler := eventHandler.NewEventHandler(eventService)

	r := router.NewRouter(guestHandler, eventHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server error:", err)
	}
}
