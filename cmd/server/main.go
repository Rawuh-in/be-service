// Package main RAWUH service
//
// @title RAWUH Service API
// @version 1.0
// @description This is the RAWUH Service API.
// @host localhost:8080
// @BasePath /
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	eventHandler "rawuh-service/internal/event/handler"
	eventDb "rawuh-service/internal/event/repository"
	eventService "rawuh-service/internal/event/service"
	guestHandler "rawuh-service/internal/guest/handler"
	guestDb "rawuh-service/internal/guest/repository"
	guestService "rawuh-service/internal/guest/service"
	projectHandler "rawuh-service/internal/project/handler"
	projectDb "rawuh-service/internal/project/repository"
	projectService "rawuh-service/internal/project/service"
	"rawuh-service/internal/shared/config"
	"rawuh-service/internal/shared/db"
	"rawuh-service/internal/shared/lib/utils"
	"rawuh-service/internal/shared/logger"
	"rawuh-service/internal/shared/redis"
	"rawuh-service/internal/shared/router"

	authHandler "rawuh-service/internal/auth/handler"
	authDb "rawuh-service/internal/auth/repository"
	authService "rawuh-service/internal/auth/service"
	userHandler "rawuh-service/internal/user/handler"
	userDb "rawuh-service/internal/user/repository"
	userService "rawuh-service/internal/user/service"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

var appConfig *config.Config
var zapLog *logger.Logger

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	appConfig = config.InitConfig()
	fluentbitPort := utils.GetEnv("FLUENTBIT_PORT", "24224")
	fluentbitPortInt, _ := strconv.Atoi(fluentbitPort)

	zapLog = logger.New(&logger.LoggerConfig{
		Env:           utils.GetEnv("ENV", "development"),
		ProductName:   "rawuh-service",
		ServiceName:   "rawuh-service",
		LogLevel:      utils.GetEnv("LOG_LEVEL", "info"),
		LogOutput:     utils.GetEnv("LOG_OUTPUT", "console"),
		FluentbitHost: utils.GetEnv("FLUENTBIT_HOST", "localhost"),
		FluentbitPort: fluentbitPortInt,
		ProcessId:     utils.GetEnv("PROCESS_ID", "rawuh-service-1"),
	})

	chosenDSN := os.Getenv("DB_DSN")
	if chosenDSN == "" {
		chosenDSN = os.Getenv("DATABASE_URL")
	}

	if chosenDSN == "" {
		log.Fatal("No database DSN found — check your environment variables")
	}

	appConfig.Dsn = chosenDSN
	zapLog.Info("Start connecting to db ", appConfig.Dsn)

	conn, err := pgx.Connect(context.Background(), chosenDSN)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer conn.Close(context.Background())

	var version string
	if err := conn.QueryRow(context.Background(), "SELECT version()").Scan(&version); err != nil {
		log.Fatalf("Query failed: %v", err)
	}

	log.Println("Connected to:", version)

	gormDB, err := config.InitDB(chosenDSN)
	if err != nil {
		zapLog.Fatal("Failed to connect to DB:", err)
	}

	zapLog.Info("Success connecting to db ", appConfig.Dsn)

	dbProvider := db.NewProvider(gormDB)
	guestDB := guestDb.NewGuestRepository(dbProvider)
	eventDB := eventDb.NewEventRepository(dbProvider)
	projectDB := projectDb.NewProjectRepository(dbProvider)
	userDB := userDb.NewUserRepository(dbProvider)

	// initialize redis (read directly from env)
	redisAddr := utils.GetEnv("REDIS_ADDR", "localhost:6379")
	redisPass := utils.GetEnv("REDIS_PASS", "")
	redisDB := 0
	rdb := redis.NewRedis(redisAddr, redisPass, redisDB)

	// repositories
	authRepo := authDb.NewAuthRepository(dbProvider)

	// services
	guestService := guestService.NewGuestService(guestDB, zapLog)
	eventService := eventService.NewEventService(eventDB, zapLog)
	userService := userService.NewUserService(userDB, authRepo, rdb, zapLog)
	projectService := projectService.NewProjectService(projectDB, zapLog)
	authService := authService.NewAuthService(authRepo, zapLog)

	// handlers
	guestHandler := guestHandler.NewGuestHandler(guestService)
	eventHandler := eventHandler.NewEventHandler(eventService)
	projectHandler := projectHandler.NewProjectHandler(projectService)
	userHandler := userHandler.NewUserHandler(userService)
	authHandler := authHandler.NewAuthHandler(authService, userDB, rdb, zapLog)

	r := router.NewRouter(guestHandler, eventHandler, projectHandler, userHandler, authHandler, rdb)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local dev
	}

	log.Println("Server running on port:", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Server error:", err)
	}

}
