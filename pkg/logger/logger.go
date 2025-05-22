package logger

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/userblog/management/pkg/config"
	"github.com/userblog/management/pkg/helper"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// contextKey is a type for context keys used by the logger
type contextKey string

// Context keys
const (
	TraceIDKey   contextKey = "trace_id"
	DebugIDKey   contextKey = "debug_id"
	ClientIpKey  contextKey = "ip"
	UserAgentKey contextKey = "user_agent"
)

var (
	log         *zap.Logger
	ProjectRoot string
)

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	return log
}

// init automatically initializes the global logger when the package is imported
func init() {
	// Initialize the project root to the working directory
	ProjectRoot, _ = os.Getwd()

	// Define the encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey, // Omit a separate function key since we include it in the caller
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   customCallerEncoder, // Use our custom caller encoder
		EncodeName:     zapcore.FullNameEncoder,
	}

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	level := zap.DebugLevel
	if config.GetOrDefaultString("ENV", "") == "production" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
		level = zap.InfoLevel
	}

	// Create a new core that writes to stdout
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(level),
	)

	// Create the logger with stack traces for errors
	log = zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel), // Add stack traces for Error level and above
	)
}

// AddToContext adds values to the context for logging
func AddToContext(ctx context.Context, key contextKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

// Debug logs a debug message
func Debug(ctx context.Context, msg string) {
	DebugF(ctx, "%s", msg)
}

// DebugF logs a formatted debug message with context values embedded in the message
func DebugF(ctx context.Context, format string, args ...interface{}) {
	log.Debug(getFormattedMsg(ctx, format, args...))
}

// Info logs an info message
func Info(ctx context.Context, msg string) {
	log.Info(getFormattedMsg(ctx, msg))
}

// InfoF logs a formatted info message with context values embedded in the message
func InfoF(ctx context.Context, format string, args ...interface{}) {
	log.Info(getFormattedMsg(ctx, format, args...))
}

// Warn logs a warning message
func Warn(ctx context.Context, msg string) {
	log.Warn(getFormattedMsg(ctx, msg))
}

// WarnF logs a formatted warning message with context values embedded in the message
func WarnF(ctx context.Context, format string, args ...interface{}) {
	log.Warn(getFormattedMsg(ctx, format, args...))
}

// Error logs an error message
func Error(ctx context.Context, msg string) {
	log.Error(getFormattedMsg(ctx, msg))
}

// ErrorF logs a formatted error message with context values embedded in the message
func ErrorF(ctx context.Context, format string, args ...interface{}) {
	log.Error(getFormattedMsg(ctx, format, args...))
}

// Fatal logs a fatal message and then exits
func Fatal(ctx context.Context, msg string) {
	log.Fatal(getFormattedMsg(ctx, msg))
}

// FatalF logs a formatted fatal message with context values embedded in the message
func FatalF(ctx context.Context, format string, args ...interface{}) {
	log.Fatal(getFormattedMsg(ctx, format, args...))
}

// getFormattedMsg creates a message with context info prefixed
func getFormattedMsg(ctx context.Context, format string, args ...interface{}) string {
	contextInfo := formatContextInfo(ctx)
	return contextInfo + fmt.Sprintf(format, args...)
}

// formatContextInfo creates a formatted string with context information
func formatContextInfo(ctx context.Context) string {
	var sb strings.Builder

	// Track if we've added any items to add appropriate separators
	hasItems := false

	// Extract goroutine ID
	goroutineID := helper.GetGoroutineID()
	if goroutineID > 0 {
		sb.WriteString("goroutine_id ")
		sb.WriteString(strconv.FormatUint(goroutineID, 10))
		hasItems = true
	}

	// Get debug ID if exists
	if debugID, ok := ctx.Value(UserAgentKey).(string); ok && debugID != "" {
		if hasItems {
			sb.WriteString(" ")
		}
		sb.WriteString("debug_id ")
		sb.WriteString(debugID)
		hasItems = true
	}

	// Get trace ID if exists
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok && traceID != "" {
		if hasItems {
			sb.WriteString(" ")
		}
		sb.WriteString("trace_id ")
		sb.WriteString(traceID)
	}

	if ip, ok := ctx.Value(ClientIpKey).(string); ok && ip != "" {
		if hasItems {
			sb.WriteString(" ")
		}
		sb.WriteString("ip ")
		sb.WriteString(ip)
	}

	sb.WriteString(" ")
	return sb.String()
}

// customCallerEncoder encodes caller information with a function name
func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// Get the file path
	filePath := caller.File
	fullFuncName := caller.Function

	if lastDotIndex := strings.LastIndex(fullFuncName, "."); lastDotIndex != -1 {
		fullFuncName = fullFuncName[lastDotIndex+1:]
	}

	// If the project root is defined, try to make the path relative
	if ProjectRoot != "" && strings.HasPrefix(filePath, ProjectRoot) {
		relPath := filePath[len(ProjectRoot)+1:]
		filePath = relPath
	}

	// Include the function name with file and line number
	enc.AppendString(fmt.Sprintf("%s:%d %s", filePath, caller.Line, fullFuncName))
}
