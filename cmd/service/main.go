package main

import (
	"context"
	"errors"
	"flappy-backend/internal/handler"
	"flappy-backend/internal/service"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pgxLogrus "github.com/jackc/pgx-logrus"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"flappy-backend/internal/db"
	"flappy-backend/pkg/logger"
)

func initializeLogger(logFile *os.File) *logrus.Logger {
	log := &logrus.Logger{
		Out:   io.MultiWriter(logFile, os.Stdout),
		Level: logrus.DebugLevel,
		Formatter: &prefixed.TextFormatter{
			DisableColors:   false,
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
			ForceFormatting: true,
		},
	}
	return log
}

func main() {
	logPath := strings.TrimRight(os.Getenv("LOG_DIR"), "/")
	if _, err := os.Stat(logPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(logPath, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
	f, err := os.OpenFile(logPath+"/"+os.Getenv("LOG_FILENAME"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create logfile " + logPath + "/" + os.Getenv("LOG_FILENAME"))
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Failed to close logfile " + logPath + "/" + os.Getenv("LOG_FILENAME"))
		}
	}(f)

	logrusLogger := initializeLogger(f)

	logrusLoggerAdapter := logger.NewLogrusAdapter(logrusLogger)

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_DATABASE"),
		os.Getenv("DB_PORT"))

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal(err)
	}

	tr := &tracelog.TraceLog{
		Logger:   pgxLogrus.NewLogger(logrusLogger),
		LogLevel: tracelog.LogLevelDebug,
	}

	config.ConnConfig.Tracer = tr
	config.MaxConns = 50
	config.MaxConnLifetime = time.Minute * 10
	config.MaxConnIdleTime = time.Minute * 3

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logrusLoggerAdapter.Fatal("failed to connect to pg db", map[string]interface{}{})
	}

	pgxDB := db.NewPgxDB(pool)

	flappyService := service.NewFlappyService(pgxDB)

	flappyHandler := handler.NewHandler(flappyService, logrusLoggerAdapter)

	app := fiber.New()
	route := app.Group("/api")
	route.Post("", flappyHandler.AddRecord)
	route.Get("", flappyHandler.GetRecord)
	route.Get("/top10", flappyHandler.GetTop10Records)

	if err = app.Listen(os.Getenv("SERVER_URL")); err != nil {
		logrusLoggerAdapter.Fatal("Server is not running!", map[string]interface{}{
			"reason": err.Error(),
		})
	}
}
