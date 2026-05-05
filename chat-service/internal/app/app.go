package app

import (
	"chat-service/internal/repository"
	httptransport "chat-service/internal/transport/http"
	"chat-service/internal/usecase"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

const serviceName = "chat-service"

func SetupChatService() (*gin.Engine, string, error) {
	port := getEnv("PORT", "8084")
	jwtSecret := getEnv("JWT_SECRET", "medical-platform-secret")
	patientServiceURL := getEnv("PATIENT_SERVICE_URL", "http://localhost:8083")
	appointmentServiceURL := getEnv("APPOINTMENT_SERVICE_URL", "http://localhost:8081")

	db, err := sql.Open("postgres", buildDatabaseURL())
	if err != nil {
		return nil, "", err
	}
	if err := db.Ping(); err != nil {
		return nil, "", err
	}
	if err := ensureSchema(db); err != nil {
		return nil, "", err
	}

	repo := repository.NewPostgresChatRepository(db)
	patientClient := usecase.NewHTTPPatientClient(patientServiceURL)
	appointmentClient := usecase.NewHTTPAppointmentClient(appointmentServiceURL)
	uc := usecase.NewChatUsecaseImpl(repo, patientClient, appointmentClient)
	handler := httptransport.NewChatHandler(uc)

	r := gin.Default()
	r.Use(httptransport.MetricsMiddleware(serviceName))
	r.GET("/metrics", httptransport.MetricsHandler)
	r.GET("/health", httptransport.HealthHandler(serviceName))

	protected := r.Group("/chat")
	protected.Use(httptransport.AuthMiddleware(jwtSecret))
	protected.POST("/threads", handler.CreateThread)
	protected.GET("/threads", handler.ListThreads)
	protected.GET("/threads/:id/messages", handler.GetMessages)
	protected.POST("/threads/:id/messages", handler.SendMessage)

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

func ensureSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_threads (
			id UUID PRIMARY KEY,
			appointment_id UUID NOT NULL UNIQUE REFERENCES appointments(id) ON DELETE CASCADE,
			patient_id UUID NOT NULL REFERENCES patient_profiles(id) ON DELETE CASCADE,
			patient_name TEXT NOT NULL,
			subject TEXT NOT NULL,
			status TEXT NOT NULL CHECK (status IN ('open', 'closed')),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			last_message_at TIMESTAMPTZ
		);

		CREATE TABLE IF NOT EXISTS chat_messages (
			id UUID PRIMARY KEY,
			thread_id UUID NOT NULL REFERENCES chat_threads(id) ON DELETE CASCADE,
			sender_user_id UUID NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
			sender_role TEXT NOT NULL CHECK (sender_role IN ('admin', 'patient')),
			sender_name TEXT NOT NULL,
			body TEXT NOT NULL CHECK (char_length(body) BETWEEN 1 AND 2000),
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_chat_threads_patient_id ON chat_threads(patient_id);
		CREATE INDEX IF NOT EXISTS idx_chat_messages_thread_id_created_at ON chat_messages(thread_id, created_at);
	`)
	return err
}
