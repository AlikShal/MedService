package http

import (
	"chat-service/internal/usecase"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const claimsContextKey = "auth_claims"

type ChatHandler struct {
	usecase usecase.ChatUsecase
}

func NewChatHandler(usecase usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{usecase: usecase}
}

func (h *ChatHandler) CreateThread(c *gin.Context) {
	var req struct {
		AppointmentID string `json:"appointment_id" binding:"required"`
		Subject       string `json:"subject"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	thread, err := h.usecase.CreateThread(req.AppointmentID, req.Subject, c.GetHeader("Authorization"), getClaims(c))
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, thread)
}

func (h *ChatHandler) ListThreads(c *gin.Context) {
	threads, err := h.usecase.ListThreads(c.GetHeader("Authorization"), getClaims(c))
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, threads)
}

func (h *ChatHandler) GetMessages(c *gin.Context) {
	messages, err := h.usecase.GetMessages(c.Param("id"), c.GetHeader("Authorization"), getClaims(c))
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req struct {
		Body string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message, err := h.usecase.SendMessage(c.Param("id"), req.Body, c.GetHeader("Authorization"), getClaims(c))
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, message)
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := usecase.ParseToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(claimsContextKey, claims)
		c.Next()
	}
}

func getClaims(c *gin.Context) *usecase.Claims {
	value, exists := c.Get(claimsContextKey)
	if !exists {
		return nil
	}
	claims, _ := value.(*usecase.Claims)
	return claims
}

func statusForError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	message := strings.ToLower(err.Error())
	switch {
	case strings.Contains(message, "missing auth"):
		return http.StatusUnauthorized
	case strings.Contains(message, "only patients"):
		return http.StatusForbidden
	case strings.Contains(message, "not found"):
		return http.StatusNotFound
	default:
		return http.StatusBadRequest
	}
}

var (
	metricsOnce sync.Once
	requests    = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "medical_platform_http_requests_total",
			Help: "Total HTTP requests handled by the service.",
		},
		[]string{"service", "method", "path", "status"},
	)
	latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "medical_platform_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "path"},
	)
)

func MetricsMiddleware(serviceName string) gin.HandlerFunc {
	metricsOnce.Do(func() {
		prometheus.MustRegister(requests, latency)
	})

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		requests.WithLabelValues(serviceName, c.Request.Method, path, http.StatusText(c.Writer.Status())).Inc()
		latency.WithLabelValues(serviceName, c.Request.Method, path).Observe(time.Since(start).Seconds())
	}
}

func MetricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func HealthHandler(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": serviceName,
		})
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || (err != nil && strings.Contains(strings.ToLower(err.Error()), "not found"))
}
