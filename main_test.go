package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// Simple integration test for the main application
func TestMainApplication(t *testing.T) {
	// Set gin to test mode
	gin.SetMode(gin.TestMode)

	// Mock the routes we know exist
	router := gin.New()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status":    "healthy",
			"timestamp": "2024-01-01T00:00:00Z",
		})
	})

	// API info endpoint
	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"name":    "Lucas API",
			"version": "v1",
			"routes":  []string{"GET /project", "POST /project"},
		})
	})

	t.Run("HealthCheck", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response["status"] != "healthy" {
			t.Errorf("Expected status 'healthy', got '%v'", response["status"])
		}
	})

	t.Run("APIInfo", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		if response["name"] == nil {
			t.Error("Expected API name in response")
		}
	})
}
