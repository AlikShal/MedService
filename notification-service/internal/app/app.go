package app

import (
	httptransport "notification-service/internal/transport/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"log/slog"
)

const defaultServiceName = "notification-service"

func SetupNotificationService() (*gin.Engine, string, error) {
	port := getEnv("PORT", "8085")
	serviceName := getEnv("SERVICE_NAME", defaultServiceName)
	notificationDelayMs := getEnvAsInt("NOTIFICATION_DELAY_MS", 150)
	logger := newLogger(getEnv("LOG_LEVEL", "info"), serviceName)
	gin.SetMode(getGinMode())

	handler := httptransport.NewNotificationHandler(serviceName, notificationDelayMs, logger)

	r := gin.New()
	r.Use(gin.Recovery())
	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, "", err
	}
	r.Use(httptransport.RequestLoggingMiddleware(logger))
	r.Use(httptransport.MetricsMiddleware(serviceName))
	r.GET("/metrics", httptransport.MetricsHandler)
	r.GET("/health", httptransport.HealthHandler(serviceName))
	r.POST("/notify", handler.Notify)

	return r, port, nil
}

func newLogger(level, serviceName string) *slog.Logger {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slogLevel,
	})).With("service", serviceName)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func getGinMode() string {
	mode := getEnv("GIN_MODE", gin.ReleaseMode)
	switch mode {
	case gin.DebugMode, gin.ReleaseMode, gin.TestMode:
		return mode
	default:
		return gin.ReleaseMode
	}
}
