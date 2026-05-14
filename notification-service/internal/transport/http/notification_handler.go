package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type NotificationHandler struct {
	serviceName string
	delayMs     int
	logger      *slog.Logger
}

type notifyRequest struct {
	Recipient string            `json:"recipient" binding:"required"`
	Channel   string            `json:"channel" binding:"required"`
	Subject   string            `json:"subject"`
	Message   string            `json:"message" binding:"required"`
	Metadata  map[string]string `json:"metadata"`
}

func NewNotificationHandler(serviceName string, delayMs int, logger *slog.Logger) *NotificationHandler {
	return &NotificationHandler{
		serviceName: serviceName,
		delayMs:     delayMs,
		logger:      logger,
	}
}

func (h *NotificationHandler) Notify(c *gin.Context) {
	var req notifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid notify request", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.delayMs > 0 {
		time.Sleep(time.Duration(h.delayMs) * time.Millisecond)
	}

	notificationID := fmt.Sprintf("notif-%d", time.Now().UnixNano())
	h.logger.Info("notification sent",
		"notification_id", notificationID,
		"recipient", req.Recipient,
		"channel", req.Channel,
		"subject", req.Subject,
		"message_length", len(req.Message),
		"metadata", req.Metadata,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":          "sent",
		"service":         h.serviceName,
		"notification_id": notificationID,
		"sent_at":         time.Now().UTC().Format(time.RFC3339),
	})
}

var (
	metricsOnce = sync.Once{}

	requestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "medical_platform_http_requests_total",
			Help: "Total HTTP requests handled by the service.",
		},
		[]string{"service", "method", "path", "status"},
	)

	errorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "medical_platform_http_errors_total",
			Help: "Total HTTP error responses from the service.",
		},
		[]string{"service", "method", "path"},
	)

	requestDuration = prometheus.NewHistogramVec(
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
		prometheus.MustRegister(requestCount, errorCount, requestDuration)
	})

	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		statusText := http.StatusText(c.Writer.Status())
		requestCount.WithLabelValues(serviceName, c.Request.Method, path, statusText).Inc()
		requestDuration.WithLabelValues(serviceName, c.Request.Method, path).Observe(time.Since(start).Seconds())
		if c.Writer.Status() >= http.StatusBadRequest {
			errorCount.WithLabelValues(serviceName, c.Request.Method, path).Inc()
		}
	}
}

func RequestLoggingMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		logger.Info("request handled",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
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
