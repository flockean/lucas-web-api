package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	API      APIConfig      `json:"api"`
	Client   ClientConfig   `json:"client"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host              string        `json:"host"`
	Port              string        `json:"port"`
	ReadTimeout       time.Duration `json:"read_timeout"`
	WriteTimeout      time.Duration `json:"write_timeout"`
	ReadHeaderTimeout time.Duration `json:"read_header_timeout"`
	ShutdownTimeout   time.Duration `json:"shutdown_timeout"`
	Debug             bool          `json:"debug"`
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	User         string        `json:"user"`
	Password     string        `json:"password"`
	Database     string        `json:"database"`
	SSLMode      string        `json:"ssl_mode"`
	MaxOpenConns int           `json:"max_open_conns"`
	MaxIdleConns int           `json:"max_idle_conns"`
	MaxLifetime  time.Duration `json:"max_lifetime"`
}

// APIConfig holds API-specific configuration
type APIConfig struct {
	Version    string        `json:"version"`
	Prefix     string        `json:"prefix"`
	Timeout    time.Duration `json:"timeout"`
	RateLimit  int           `json:"rate_limit"`
	EnableCORS bool          `json:"enable_cors"`
	EnableAuth bool          `json:"enable_auth"`
	JWTSecret  string        `json:"jwt_secret"`
}

// ClientConfig holds client-specific configuration
type ClientConfig struct {
	BaseURL       string        `json:"base_url"`
	Timeout       time.Duration `json:"timeout"`
	RetryCount    int           `json:"retry_count"`
	RetryWaitTime time.Duration `json:"retry_wait_time"`
	Debug         bool          `json:"debug"`
	UserAgent     string        `json:"user_agent"`
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Host:              getEnvOrDefault("SERVER_HOST", "0.0.0.0"),
			Port:              getEnvOrDefault("SERVER_PORT", "8080"),
			ReadTimeout:       getDurationEnvOrDefault("SERVER_READ_TIMEOUT", 30*time.Second),
			WriteTimeout:      getDurationEnvOrDefault("SERVER_WRITE_TIMEOUT", 30*time.Second),
			ReadHeaderTimeout: getDurationEnvOrDefault("SERVER_READ_HEADER_TIMEOUT", 5*time.Second),
			ShutdownTimeout:   getDurationEnvOrDefault("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second),
			Debug:             getBoolEnvOrDefault("SERVER_DEBUG", false),
		},
		Database: DatabaseConfig{
			Host:         getEnvOrDefault("DB_HOST", "0.0.0.0"),
			Port:         getIntEnvOrDefault("DB_PORT", 5432),
			User:         getEnvOrDefault("DB_USER", "postgres"),
			Password:     getEnvOrDefault("DB_PASSWORD", "rocket123"),
			Database:     getEnvOrDefault("DB_NAME", "projectDb"),
			SSLMode:      getEnvOrDefault("DB_SSL_MODE", "disable"),
			MaxOpenConns: getIntEnvOrDefault("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getIntEnvOrDefault("DB_MAX_IDLE_CONNS", 10),
			MaxLifetime:  getDurationEnvOrDefault("DB_MAX_LIFETIME", 5*time.Minute),
		},
		API: APIConfig{
			Version:    getEnvOrDefault("API_VERSION", "v1"),
			Prefix:     getEnvOrDefault("API_PREFIX", "/api"),
			Timeout:    getDurationEnvOrDefault("API_TIMEOUT", 30*time.Second),
			RateLimit:  getIntEnvOrDefault("API_RATE_LIMIT", 100),
			EnableCORS: getBoolEnvOrDefault("API_ENABLE_CORS", true),
			EnableAuth: getBoolEnvOrDefault("API_ENABLE_AUTH", false),
			JWTSecret:  getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		},
		Client: ClientConfig{
			BaseURL:       getEnvOrDefault("CLIENT_BASE_URL", "http://localhost:8080/api"),
			Timeout:       getDurationEnvOrDefault("CLIENT_TIMEOUT", 30*time.Second),
			RetryCount:    getIntEnvOrDefault("CLIENT_RETRY_COUNT", 3),
			RetryWaitTime: getDurationEnvOrDefault("CLIENT_RETRY_WAIT_TIME", 2*time.Second),
			Debug:         getBoolEnvOrDefault("CLIENT_DEBUG", false),
			UserAgent:     getEnvOrDefault("CLIENT_USER_AGENT", "Lucas-API-Client/1.0"),
		},
	}
}

// GetConnectionString returns the database connection string
func (c *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

// GetServerAddress returns the server address
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// Helper functions
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnvOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnvOrDefault(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnvOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}

	if c.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	return nil
}
