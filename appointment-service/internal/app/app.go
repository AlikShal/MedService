package app

import (
	"appointment-service/internal/repository"
	httptransport "appointment-service/internal/transport/http"
	"appointment-service/internal/usecase"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

const serviceName = "appointment-service"

func SetupAppointmentService() (*gin.Engine, string, error) {
	port := getEnv("PORT", "8081")
	jwtSecret := getEnv("JWT_SECRET", "medical-platform-secret")
	doctorServiceURL := getEnv("DOCTOR_SERVICE_URL", "http://localhost:8080")
	patientServiceURL := getEnv("PATIENT_SERVICE_URL", "http://localhost:8083")

	db, err := sql.Open("postgres", buildDatabaseURL())
	if err != nil {
		return nil, "", err
	}
	if err := db.Ping(); err != nil {
		return nil, "", err
	}

	repo := repository.NewPostgresAppointmentRepository(db)
	doctorClient := usecase.NewHTTPDoctorClient(doctorServiceURL)
	patientClient := usecase.NewHTTPPatientClient(patientServiceURL)
	uc := usecase.NewAppointmentUsecaseImpl(repo, doctorClient, patientClient)
	handler := httptransport.NewAppointmentHandler(uc)

	gin.SetMode(getEnv("GIN_MODE", gin.ReleaseMode))
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, "", err
	}
	r.Use(httptransport.MetricsMiddleware(serviceName))
	r.GET("/metrics", httptransport.MetricsHandler)
	r.GET("/health", httptransport.HealthHandler(serviceName))

	protected := r.Group("/appointments")
	protected.Use(httptransport.AuthMiddleware(jwtSecret))
	protected.POST("", httptransport.RequireRoles("patient"), handler.CreateAppointment)
	protected.GET("/:id", handler.GetAppointmentByID)
	protected.GET("", handler.GetAllAppointments)
	protected.PATCH("/:id/status", httptransport.RequireRoles("admin"), handler.UpdateAppointmentStatus)

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
