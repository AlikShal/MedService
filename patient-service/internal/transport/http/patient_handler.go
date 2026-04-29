package http

import (
	"database/sql"
	"errors"
	"net/http"
	"patient-service/internal/usecase"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const claimsContextKey = "auth_claims"

type PatientHandler struct {
	usecase usecase.PatientUsecase
}

func NewPatientHandler(usecase usecase.PatientUsecase) *PatientHandler {
	return &PatientHandler{usecase: usecase}
}

func (h *PatientHandler) SaveProfile(c *gin.Context) {
	claims := getClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	var req struct {
		FullName    string `json:"full_name" binding:"required"`
		Phone       string `json:"phone"`
		DateOfBirth string `json:"date_of_birth"`
		Notes       string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patient, err := h.usecase.SaveProfile(claims.UserID, claims.Email, req.FullName, req.Phone, req.DateOfBirth, req.Notes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) GetMyProfile(c *gin.Context) {
	claims := getClaims(c)
	if claims == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	patient, err := h.usecase.GetByUserID(claims.UserID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) GetByID(c *gin.Context) {
	patient, err := h.usecase.GetByID(c.Param("id"))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) GetAll(c *gin.Context) {
	patients, err := h.usecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, patients)
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

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := getClaims(c)
		if claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
			return
		}

		if !usecase.HasRole(claims.Role, roles...) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

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

func HealthHandler(db *sql.DB, serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"service": serviceName,
				"status":  "degraded",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"status":  "ok",
			"storage": "postgres",
		})
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || (err != nil && strings.Contains(strings.ToLower(err.Error()), "not found"))
}
