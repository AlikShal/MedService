package http

import (
	"database/sql"
	"doctor-service/internal/usecase"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type DoctorHandler struct {
	usecase usecase.DoctorUsecase
}

func NewDoctorHandler(usecase usecase.DoctorUsecase) *DoctorHandler {
	return &DoctorHandler{usecase: usecase}
}

func (h *DoctorHandler) CreateDoctor(c *gin.Context) {
	var req struct {
		FullName       string `json:"full_name" binding:"required"`
		Specialization string `json:"specialization"`
		Email          string `json:"email" binding:"required"`
		Office         string `json:"office"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	doctor, err := h.usecase.CreateDoctor(req.FullName, req.Specialization, req.Email, req.Office)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, doctor)
}

func (h *DoctorHandler) GetDoctorByID(c *gin.Context) {
	id := c.Param("id")
	doctor, err := h.usecase.GetDoctorByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctor)
}

func (h *DoctorHandler) GetAllDoctors(c *gin.Context) {
	doctors, err := h.usecase.GetAllDoctors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, doctors)
}

const claimsContextKey = "auth_claims"

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(claimsContextKey, claims)
		c.Next()
	}
}

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		value, exists := c.Get(claimsContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
			return
		}

		claims := value.(*Claims)
		for _, role := range roles {
			if strings.EqualFold(claims.Role, role) {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
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
