package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/userblog/management/pkg/logger"
	"net/http"
)

// ErrorResponse is the standard error response format
type ErrorResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId,omitempty"`
}

// GlobalExceptionHandler catches all unhandled panics and returns a proper error response
func GlobalExceptionHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := c.Request.Context()
				traceID, _ := ctx.Value(logger.TraceIDKey).(string)
				logger.ErrorF(ctx, "PANIC RECOVERED: %v", err)

				// Determine if response was already written
				if c.Writer.Written() {
					return
				}

				// Create a consistent error response
				errorResponse := ErrorResponse{
					Status:  http.StatusInternalServerError,
					Code:    "INTERNAL_SERVER_ERROR",
					Message: "An unexpected error occurred. Our team has been notified.",
					TraceID: traceID,
				}

				// Abort the request with an internal server error
				c.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse)
			}
		}()

		// Continue to next handler
		c.Next()

		// Check for any errors added during request handling
		if len(c.Errors) > 0 {
			// Get the current context with all the trace information
			ctx := c.Request.Context()

			// Get trace ID from context
			traceID, _ := ctx.Value(logger.TraceIDKey).(string)

			// Log all errors
			for _, e := range c.Errors {
				logger.Error(ctx, fmt.Sprintf("Request Error: %v", e.Err))
			}

			// If no response was sent yet, send a standard error response
			if !c.Writer.Written() {
				// Get the last error to determine response
				lastError := c.Errors.Last()

				// Default values
				status := http.StatusInternalServerError
				code := "INTERNAL_SERVER_ERROR"
				message := "An unexpected error occurred"

				// Check if the error has status code information
				if lastError.Type == gin.ErrorTypePublic {
					// For public errors, use the error message directly
					message = lastError.Error()
				}

				// Check for specific error types
				switch typed := lastError.Err.(type) {
				case *gin.Error:
					// Use the gin error's status code if available
					if typed.Meta != nil {
						if code, ok := typed.Meta.(int); ok {
							status = code
						}
					}
				case interface{ StatusCode() int }:
					// For errors that implement StatusCode method
					status = typed.StatusCode()
				case interface{ Code() string }:
					// For errors that implement Code method
					code = typed.Code()
				}

				// Map HTTP status codes to error codes if needed
				if status == http.StatusNotFound {
					code = "NOT_FOUND"
					message = "The requested resource was not found"
				} else if status == http.StatusUnauthorized {
					code = "UNAUTHORIZED"
					message = "Authentication is required to access this resource"
				} else if status == http.StatusForbidden {
					code = "FORBIDDEN"
					message = "You don't have permission to access this resource"
				} else if status == http.StatusBadRequest {
					code = "BAD_REQUEST"
					if message == "An unexpected error occurred" {
						message = "The request was invalid or cannot be served"
					}
				}

				// Create the error response
				errorResponse := ErrorResponse{
					Status:  status,
					Code:    code,
					Message: message,
					TraceID: traceID,
				}

				// Respond with the appropriate status code and error message
				c.AbortWithStatusJSON(status, errorResponse)
			}
		}
	}
}
