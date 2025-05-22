package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/userblog/management/pkg/logger"
)

// Logger middleware adds request logging and trace ID to context
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {

		start := time.Now()
		ctx := c.Request.Context()
		clientIP := c.ClientIP()
		traceID := uuid.New().String()
		userAgent := c.Request.UserAgent()

		c.Request = c.Request.WithContext(ctx)
		c.Header("X-Trace-ID", traceID)

		method := c.Request.Method
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Add query string if present
		fullPath := path
		if raw != "" {
			fullPath = path + "?" + raw
		}

		c.Set(string(logger.ClientIpKey), clientIP)
		c.Set(string(logger.UserAgentKey), userAgent)
		c.Set(string(logger.TraceIDKey), traceID)
		c.Next()

		// Calculate request duration
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()

		logger.InfoF(ctx, "API RESPONSE "+
			"Method: %s | "+
			"Path: %s | "+
			"Status: %d | "+
			"Size: %d bytes | "+
			"Duration: %v",
			method,
			fullPath,
			statusCode,
			responseSize,
			duration)

	}
}
