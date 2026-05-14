package app

import (
	"database/sql"
	"doctor-service/internal/repository"
	httptransport "doctor-service/internal/transport/http"
	"doctor-service/internal/usecase"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

const serviceName = "doctor-service"

func SetupDoctorService() (*gin.Engine, string, error) {
	port := getEnv("PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "medical-platform-secret")

	db, err := sql.Open("postgres", buildDatabaseURL())
	if err != nil {
		return nil, "", err
	}
	if err := db.Ping(); err != nil {
		return nil, "", err
	}

	repo := repository.NewPostgresDoctorRepository(db)
	uc := usecase.NewDoctorUsecaseImpl(repo)
	handler := httptransport.NewDoctorHandler(uc)

	gin.SetMode(getEnv("GIN_MODE", gin.ReleaseMode))
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, "", err
	}
	r.Use(httptransport.MetricsMiddleware(serviceName))
	r.GET("/metrics", httptransport.MetricsHandler)
	r.GET("/health", httptransport.HealthHandler(serviceName))
	r.GET("/doctors/:id", handler.GetDoctorByID)
	r.GET("/doctors", handler.GetAllDoctors)

	admin := r.Group("/doctors")
	admin.Use(httptransport.AuthMiddleware(jwtSecret), httptransport.RequireRoles("admin"))
	admin.POST("", handler.CreateDoctor)

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
