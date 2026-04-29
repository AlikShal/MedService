package http

import (
	"auth-service/internal/model"
	"auth-service/internal/usecase"
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

type AuthHandler struct {
	usecase usecase.AuthUsecase
}

func NewAuthHandler(usecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.usecase.Register(req.FullName, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.usecase.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) Me(c *gin.Context) {
	claimsValue, exists := c.Get(claimsContextKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
		return
	}

	claims := claimsValue.(*usecase.Claims)
	user, err := h.usecase.GetUserByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.Sanitize()})
}

const claimsContextKey = "auth_claims"

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

func RequireRoles(allowed ...model.Role) gin.HandlerFunc {
	allowedRoles := make(map[model.Role]struct{}, len(allowed))
	for _, role := range allowed {
		allowedRoles[role] = struct{}{}
	}

	return func(c *gin.Context) {
		value, exists := c.Get(claimsContextKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth context"})
			return
		}

		claims := value.(*usecase.Claims)
		if _, ok := allowedRoles[claims.Role]; !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows) || (err != nil && strings.Contains(strings.ToLower(err.Error()), "not found"))
}
