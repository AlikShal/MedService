package app

import (
	"auth-service/internal/model"
	"auth-service/internal/repository"
	httptransport "auth-service/internal/transport/http"
	"auth-service/internal/usecase"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

const serviceName = "auth-service"

func SetupAuthService() (*gin.Engine, string, error) {
	port := getEnv("PORT", "8082")
	jwtSecret := getEnv("JWT_SECRET", "medical-platform-secret")
	tokenTTL := getEnvDuration("JWT_TTL", 24*time.Hour)

	db, err := sql.Open("postgres", buildDatabaseURL())
	if err != nil {
		return nil, "", err
	}
	if err := db.Ping(); err != nil {
		return nil, "", err
	}

	repo := repository.NewPostgresAuthRepository(db)
	if err := ensureAdminUser(repo); err != nil {
		return nil, "", err
	}
	uc := usecase.NewAuthUsecase(repo, jwtSecret, tokenTTL)
	handler := httptransport.NewAuthHandler(uc)

	r := gin.Default()
	r.Use(httptransport.MetricsMiddleware(serviceName))
	r.GET("/metrics", httptransport.MetricsHandler)
	r.GET("/health", httptransport.HealthHandler(serviceName))
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.Login)

	authenticated := r.Group("/auth")
	authenticated.Use(httptransport.AuthMiddleware(jwtSecret))
	authenticated.GET("/me", handler.Me)

	return r, port, nil
}

func ensureAdminUser(repo repository.AuthRepository) error {
	adminEmail := getEnv("ADMIN_EMAIL", "admin@medsync.local")
	if _, err := repo.GetUserByEmail(adminEmail); err == nil {
		return nil
	} else if !strings.Contains(strings.ToLower(err.Error()), "not found") {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(getEnv("ADMIN_PASSWORD", "admin123")), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return repo.CreateUser(&model.User{
		ID:           uuid.NewString(),
		FullName:     getEnv("ADMIN_FULL_NAME", "Clinic Administrator"),
		Email:        adminEmail,
		PasswordHash: string(passwordHash),
		Role:         model.RoleAdmin,
	})
}

func buildDatabaseURL() string {
	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		return databaseURL
	}

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	name := getEnv("DB_NAME", "medical_platform")
	sslMode := getEnv("DB_SSLMODE", "disable")

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host,
		port,
		user,
		password,
		name,
		sslMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		parsed, err := time.ParseDuration(value)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
