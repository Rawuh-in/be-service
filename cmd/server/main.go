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
	"rawuh-service/internal/shared/router"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
)

var appConfig *config.Config
var zapLog *logger.Logger

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
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
		chosenDSN = appConfig.Dsn
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

	// rdb := redis.NewRedis(appConfig.RedisAddr, appConfig.RedisPass, appConfig.RedisDB)

	// guestService := guestService.NewGuestService(guestDB, zapLog, rdb)
	guestService := guestService.NewGuestService(guestDB, zapLog)
	guestHandler := guestHandler.NewGuestHandler(guestService)

	eventService := eventService.NewEventService(eventDB, zapLog)
	eventHandler := eventHandler.NewEventHandler(eventService)

	projectService := projectService.NewProjectService(projectDB, zapLog)
	projectHandler := projectHandler.NewProjectHandler(projectService)

	r := router.NewRouter(guestHandler, eventHandler, projectHandler)

	log.Println("Server running on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal("Server error:", err)
	}
}
