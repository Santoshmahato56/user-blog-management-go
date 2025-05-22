package helper

import (
	"runtime"
	"strconv"
	"strings"
)

// GetGoroutineID extracts the current goroutine ID from runtime stack
func GetGoroutineID() uint64 {
	// This is a hack but it's the only reliable way to get goroutine ID
	buf := make([]byte, 64)
	n := runtime.Stack(buf, false)

	// Goroutine ID is the first number after "goroutine " in the stack trace
	// Format is "goroutine 42 [running]:"
	idField := strings.Fields(string(buf[:n]))[1]
	id, _ := strconv.ParseUint(idField, 10, 64)
	return id
}
