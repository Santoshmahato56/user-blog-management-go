package config

import (
	"sync"
)

// Test helpers
var testEnvironmentSetup sync.Once

// SetupTestEnvironment initializes a test environment by replacing
// the panic with a warning message and setting required test variables
func SetupTestEnvironment() {
	testEnvironmentSetup.Do(func() {
		// Modify the package-level variables
		configCache = make(map[string]interface{})
		envCache = make(map[string]string)

		// Set test values
		envCache["JWT_SECRET"] = "test-secret-key"
		envCache["TOKEN_EXPIRY"] = "24"

		// Mark cache as initialized to prevent further issues
		cacheInitialized = true
	})
}
