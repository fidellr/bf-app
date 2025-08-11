// @title Book Management API
// @version 1.0
// @description This is a book management server for a RESTful API.

// @contact.name Fidel Ramadhan
// @contact.email fidelramadhan@gmail.com

// @host localhost:8080
// @BasePath /api
// @schemes http

package main

import (
	_ "bf-api/docs" // Required for Swagger
	"bf-api/internal/app/handlers"
	"bf-api/internal/app/routes"
	"bf-api/internal/domain/services"
	"bf-api/internal/infrastructure/db/postgres"
	"bf-api/internal/infrastructure/logger"
	"bufio"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func main() {
	logger.Init(false)
	zap.ReplaceGlobals(logger.Logger)
	defer logger.Logger.Sync()

	cfg := postgres.DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "bookdb"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pgPool, err := postgres.NewPostgresDB(ctx, cfg)
	if err != nil {
		logger.Logger.Fatal("failed to connect to database", zap.Error(err))

	}
	defer pgPool.Close()

	if err := postgres.HealthCheck(ctx, pgPool); err != nil {
		log.Fatalf("Database health check failed: %v", err)
	}

	bookRepo := postgres.NewBookRepository(pgPool)
	bookSvc := services.NewBookService(bookRepo)

	e := echo.New()
	e.HideBanner = true

	bookHandler := handlers.NewBookHandler(bookSvc, logger.Logger)
	routes.APIRouter(e, bookHandler, bookSvc, logger.Logger)
	startServer(e)

}

func startServer(e *echo.Echo) {
	go func() {
		port := getEnv("PORT", "8080")
		logger.Logger.Info("Starting server", zap.String("port", port))
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			logger.Logger.Fatal("shutting down the server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logger.Logger.Info("Shutting down server...")
	if err := e.Shutdown(ctx); err != nil {
		logger.Logger.Error("Server shutdown failed", zap.Error(err))
	}
}

func getEnv(key, defaultValue string) string {
	if err := loadEnvFile(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func loadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue // Skip malformed lines
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		os.Setenv(key, value) // Set in system environment
	}
	return scanner.Err()
}
