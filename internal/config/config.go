// internal/config/config.go
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	DatabaseURL      string
	ServerPort       string
	Environment      string
	LogLevel         string
	ConnectionMaxAge int
	MaxOpenConns     int
	MaxIdleConns     int
}

// Load loads configuration from environment variables with sensible defaults
func Load() *Config {
	// Default port is 9090 as required by the problem statement
	port := getEnv("PORT", "9090")

	// Database configuration
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "order_service")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)

	// Connection pool settings
	maxOpenConns, _ := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _ := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "5"))
	connMaxAge, _ := strconv.Atoi(getEnv("DB_CONN_MAX_AGE", "300"))

	return &Config{
		DatabaseURL:      dbURL,
		ServerPort:       port,
		Environment:      getEnv("ENV", "development"),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
		ConnectionMaxAge: connMaxAge,
		MaxOpenConns:     maxOpenConns,
		MaxIdleConns:     maxIdleConns,
	}
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
