package app

import (
	"database/sql"
	"fmt"
	"os"
	"patient-service/internal/repository"
	httptransport "patient-service/internal/transport/http"
	"patient-service/internal/usecase"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

const serviceName = "patient-service"

func SetupPatientService() (*gin.Engine, string, error) {
	port := getEnv("PORT", "8083")
	jwtSecret := getEnv("JWT_SECRET", "medical-platform-secret")

	db, err := sql.Open("postgres", buildDatabaseURL())
	if err != nil {
		return nil, "", err
	}
	if err := db.Ping(); err != nil {
		return nil, "", err
	}

	repo := repository.NewPostgresPatientRepository(db)
	uc := usecase.NewPatientUsecase(repo)
	handler := httptransport.NewPatientHandler(uc)

	r := gin.Default()
	r.Use(httptransport.MetricsMiddleware(serviceName))
	r.GET("/metrics", httptransport.MetricsHandler)
	r.GET("/health", httptransport.HealthHandler(serviceName))

	protected := r.Group("/patients")
	protected.Use(httptransport.AuthMiddleware(jwtSecret))
	protected.PUT("/me", handler.SaveProfile)
	protected.GET("/me", handler.GetMyProfile)
	protected.GET("/:id", httptransport.RequireRoles("admin"), handler.GetByID)
	protected.GET("", httptransport.RequireRoles("admin"), handler.GetAll)

	return r, port, nil
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
