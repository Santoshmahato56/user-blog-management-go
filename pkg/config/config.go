package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	// Cache for all configuration values
	configCache      map[string]interface{} // Holds all config values in a flattened structure
	envCache         map[string]string      // Holds environment variables
	cacheMutex       sync.RWMutex
	cacheInitialized bool
	configPath       = "config.yml" // Default config file path
)

func init() {
	// Load environment variables
	//err := godotenv.Load()
	//if err != nil {
	//	panic("Error loading .env file, using default values")
	//}
	initCache()
}

// SetConfigPath sets a custom path for the config file
func SetConfigPath(path string) {
	configPath = path
	// Reset cache so it will be reloaded
	resetCache()
}

// resetCache clears the configuration cache
func resetCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	configCache = nil
	envCache = nil
	cacheInitialized = false
}

// ensureCacheInitialized makes sure the cache is loaded
func ensureCacheInitialized() {
	cacheMutex.RLock()
	initialized := cacheInitialized
	cacheMutex.RUnlock()

	if !initialized {
		initCache()
	}
}

// initCache initializes the configuration cache
func initCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	if cacheInitialized {
		return // Double-check inside lock
	}

	// Initialize caches
	configCache = make(map[string]interface{})
	envCache = make(map[string]string)

	// Load all environment variables into envCache
	for _, envVar := range os.Environ() {
		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) == 2 {
			envCache[pair[0]] = pair[1]
		}
	}

	// Load config file into configCache
	loadConfigFile()

	cacheInitialized = true
}

// flattenMap recursively flattens a nested map into a single-level map with key paths
func flattenMap(result map[string]interface{}, current map[string]interface{}, prefix string) {
	for k, v := range current {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch child := v.(type) {
		case map[string]interface{}:
			// Recursively flatten nested maps
			flattenMap(result, child, key)
		case map[interface{}]interface{}:
			// Convert keys to strings and flatten
			stringMap := make(map[string]interface{})
			for mk, mv := range child {
				if strKey, ok := mk.(string); ok {
					stringMap[strKey] = mv
				}
			}
			flattenMap(result, stringMap, key)
		default:
			// Store the value with the full key path
			result[key] = v

			// Also store with uppercase key for easier matching with env vars
			result[strings.ToUpper(key)] = v

			// Store with underscores instead of dots for ENV_VAR style matching
			underscoreKey := strings.ReplaceAll(strings.ToUpper(key), ".", "_")
			result[underscoreKey] = v
		}
	}
}

// loadConfigFile loads the configuration from config.yml
func loadConfigFile() {
	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		// File doesn't exist or can't be read
		return
	}

	// Parse YAML into a generic map
	var rawConfig map[string]interface{}
	err = yaml.Unmarshal(data, &rawConfig)
	if err != nil {
		fmt.Printf("Warning: Error parsing config file %s: %v\n", configPath, err)
		return
	}

	// Flatten the nested structure
	flattenMap(configCache, rawConfig, "")
}

// Get retrieves a configuration value by key, returns nil if not found
func Get(key string) interface{} {
	// Ensure cache is initialized
	ensureCacheInitialized()

	// First priority: Check environment variables
	cacheMutex.RLock()
	envValue, envExists := envCache[key]
	cacheMutex.RUnlock()

	if envExists && envValue != "" {
		return envValue
	}

	// Second priority: Check config file values
	cacheMutex.RLock()
	configValue, configExists := configCache[key]
	cacheMutex.RUnlock()

	if configExists && configValue != nil {
		return configValue
	}

	// Try different key formats
	keys := []string{
		strings.ToUpper(key),
		strings.ReplaceAll(key, "_", "."),
		strings.ReplaceAll(strings.ToUpper(key), "_", "."),
	}

	for _, altKey := range keys {
		if altKey == key {
			continue // Skip if it's the same as the original key
		}

		cacheMutex.RLock()
		configValue, configExists = configCache[altKey]
		cacheMutex.RUnlock()

		if configExists && configValue != nil {
			return configValue
		}
	}

	return nil
}

// GetOrDefaultString gets a value as a string, with a fallback default
func GetOrDefaultString(key, defaultValue string) string {
	value := Get(key)
	if value == nil {
		return defaultValue
	}

	switch v := value.(type) {
	case string:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		return fmt.Sprintf("%v", v)
	default:
		return defaultValue
	}
}

// GetOrDefaultInt gets a value as an int, with a fallback default
func GetOrDefaultInt(key string, defaultValue int) int {
	value := Get(key)
	if value == nil {
		return defaultValue
	}

	switch v := value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
		return defaultValue
	case bool:
		if v {
			return 1
		}
		return 0
	default:
		return defaultValue
	}
}

// GetOrDefaultBool gets a value as a bool, with a fallback default
func GetOrDefaultBool(key string, defaultValue bool) bool {
	value := Get(key)
	if value == nil {
		return defaultValue
	}

	switch v := value.(type) {
	case bool:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v != 0
	case float32, float64:
		return v != 0
	case string:
		lower := strings.ToLower(v)
		if lower == "true" || lower == "yes" || lower == "1" || lower == "on" || lower == "enabled" {
			return true
		}
		if lower == "false" || lower == "no" || lower == "0" || lower == "off" || lower == "disabled" {
			return false
		}
		return defaultValue
	default:
		return defaultValue
	}
}

// GetOrDefaultFloat gets a value as a float64, with a fallback default
func GetOrDefaultFloat(key string, defaultValue float64) float64 {
	value := Get(key)
	if value == nil {
		return defaultValue
	}

	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f
		}
		return defaultValue
	case bool:
		if v {
			return 1.0
		}
		return 0.0
	default:
		return defaultValue
	}
}

// GetOrDefaultStringSlice gets a value as a string slice, with a fallback default
func GetOrDefaultStringSlice(key string, defaultValue []string) []string {
	value := Get(key)
	if value == nil {
		return defaultValue
	}

	switch v := value.(type) {
	case []string:
		return v
	case string:
		// Split by comma
		if v == "" {
			return []string{}
		}
		return strings.Split(v, ",")
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			} else {
				result = append(result, fmt.Sprintf("%v", item))
			}
		}
		return result
	default:
		return defaultValue
	}
}
