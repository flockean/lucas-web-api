package clients_test

import (
	"LucasApi/api/clients"
	"LucasApi/api/config"
	"LucasApi/api/dbutils"
	"testing"
	"time"
)

func TestClientGenerator(t *testing.T) {
	// Test basic client generation
	generator := clients.NewClientGenerator(
		clients.WithBaseURL("http://localhost:8080/api"),
		clients.WithTimeout(10*time.Second),
		clients.WithDebug(true),
		clients.WithMiddleware(clients.AuthMiddleware("test-token")),
		clients.WithMiddleware(clients.JSONResponseMiddleware()),
		clients.WithMiddleware(clients.RateLimitMiddleware(100)),
	)

	if generator == nil {
		t.Fatal("Expected client generator to be created")
	}

	// Test client creation
	client := generator.GenerateClient("TestClient")
	if client == nil {
		t.Fatal("Expected client to be created")
	}

	// Test generic API client
	apiClient := clients.NewGenericAPIClient(generator, "/projects")
	if apiClient == nil {
		t.Fatal("Expected generic API client to be created")
	}
}

func TestProjectsClient(t *testing.T) {
	cfg := config.LoadConfig()

	client := clients.NewProjectsClient(&cfg.Client)
	if client == nil {
		t.Fatal("Expected projects client to be created")
	}

	// Note: These tests would require a running server
	// In a real scenario, you would mock the HTTP responses
}

func ExampleClientGenerator() {
	// Create a new client generator with custom configuration
	generator := clients.NewClientGenerator(
		clients.WithBaseURL("https://api.example.com"),
		clients.WithTimeout(30*time.Second),
		clients.WithDebug(false),
		clients.WithMiddleware(clients.AuthMiddleware("your-api-token")),
		clients.WithMiddleware(clients.JSONResponseMiddleware()),
	)

	// Load configuration
	cfg := config.LoadConfig()

	// Generate a projects client
	projectsClient := clients.NewProjectsClient(&cfg.Client)

	// Create a new project
	newProject := dbutils.CreateProjectRequest{
		Name:        "My New Project",
		Description: stringPtr("A project created via REST client"),
		Status:      "active",
	}

	project, err := projectsClient.CreateProject(newProject)
	if err != nil {
		// Handle error
		return
	}

	// Use the created project
	_ = project
	_ = generator // Use generator to avoid unused variable error
}

// Helper function for example
func stringPtr(s string) *string {
	return &s
}
