package config

import (
	"log"
	"os"
	"strconv"
)

// AppConfig structure for environment-based configurations.
type AppConfig struct {
	Debug       bool `json:"debug"`
	Port        int  `json:"port"`
	MetricsPort int  `json:"metricsPort"`
}

// CFG is the global configuration object.
var CFG AppConfig

// LoadConfiguration loads configuration from environment variables.
func LoadConfiguration() {
	CFG.Debug = parseEnvBool("DEBUG", false)            // Assuming false as the default value
	CFG.Port = parseEnvInt("PORT", 8080)                // Assuming 8080 as the default port
	CFG.MetricsPort = parseEnvInt("METRICS_PORT", 9090) // Assuming 9090 as the default port
}

// getEnvOrDefault returns the value of an environment variable or a default value.
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// parseEnvInt parses an environment variable as an integer or returns a default value.
func parseEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error parsing %s as int: %v. Using default value: %d", key, err, defaultValue)
		return defaultValue
	}
	return intValue
}

// parseEnvBool parses an environment variable as a boolean or returns a default value.
func parseEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Error parsing %s as bool: %v. Using default value: %t", key, err, defaultValue)
		return defaultValue
	}
	return boolValue
}
