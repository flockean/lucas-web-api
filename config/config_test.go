package config_test

import (
	"LucasApi/api/config"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Run("LoadConfigWithDefaults", func(t *testing.T) {
		cfg := config.LoadConfig()

		// Test that required fields are not empty
		if cfg.Server.Host == "" {
			t.Error("Expected server host to be set")
		}
		if cfg.Server.Port == "" {
			t.Error("Expected server port to be set")
		}
		if cfg.Database.Host == "" {
			t.Error("Expected database host to be set")
		}
		if cfg.API.Prefix == "" {
			t.Error("Expected API prefix to be set")
		}
	})

	t.Run("LoadConfigWithEnvironmentVars", func(t *testing.T) {
		// Set custom env vars
		originalHost := os.Getenv("SERVER_HOST")
		originalPort := os.Getenv("SERVER_PORT")

		_ = os.Setenv("SERVER_HOST", "test.example.com")
		_ = os.Setenv("SERVER_PORT", "9000")

		defer func() {
			if originalHost == "" {
				_ = os.Unsetenv("SERVER_HOST")
			} else {
				_ = os.Setenv("SERVER_HOST", originalHost)
			}
			if originalPort == "" {
				_ = os.Unsetenv("SERVER_PORT")
			} else {
				_ = os.Setenv("SERVER_PORT", originalPort)
			}
		}()

		cfg := config.LoadConfig()

		// Test that environment variables are used
		if cfg.Server.Host != "test.example.com" {
			t.Errorf("Expected host 'test.example.com', got '%s'", cfg.Server.Host)
		}
		if cfg.Server.Port != "9000" {
			t.Errorf("Expected port '9000', got '%s'", cfg.Server.Port)
		}
	})
}

func TestConfigValidation(t *testing.T) {
	t.Run("ValidConfig", func(t *testing.T) {
		cfg := config.LoadConfig()

		// Test validation doesn't return error for default config
		if err := cfg.Validate(); err != nil {
			t.Errorf("Expected valid default config, got error: %v", err)
		}
	})
}

func TestServerConfig(t *testing.T) {
	t.Run("GetServerAddress", func(t *testing.T) {
		cfg := &config.ServerConfig{
			Host: "localhost",
			Port: "8080",
		}

		expected := "localhost:8080"
		if addr := cfg.GetServerAddress(); addr != expected {
			t.Errorf("Expected address '%s', got '%s'", expected, addr)
		}
	})
}
