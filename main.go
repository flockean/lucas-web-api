package main

import (
	"LucasApi/api/config"
	"LucasApi/api/controllers"
	"LucasApi/api/database"
	_ "LucasApi/api/docs"
	"LucasApi/api/services"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Lucas Web API
// @version 1.0
// @description REST API for project and service management with clean architecture
// @termsOfService http://localhost:8080/terms/

// @contact.name API Support
// @contact.url http://localhost:8080/support
// @contact.email support@lucas-api.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api
func main() {
	// Load configuration
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Initialize database
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Error closing database: %v", closeErr)
		}
	}()

	// Initialize repositories
	projectRepo := database.NewProjectRepository(db)
	serviceRepo := database.NewServiceRepository(db)

	// Initialize services
	projectService := services.NewProjectService(projectRepo, serviceRepo)
	serviceService := services.NewServiceService(serviceRepo, projectRepo)

	// Initialize controllers
	projectController := controllers.NewProjectController(projectService)
	serviceController := controllers.NewServiceController(serviceService)
	healthController := controllers.NewHealthController()
	apiController := controllers.NewAPIController()

	// Setup router
	router := setupRouter(cfg, projectController, serviceController, healthController, apiController)

	// Create server
	srv := &http.Server{
		Addr:              cfg.Server.GetServerAddress(),
		Handler:           router,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(cfg *config.Config,
	projectController *controllers.ProjectController,
	serviceController *controllers.ServiceController,
	healthController *controllers.HealthController,
	apiController *controllers.APIController) *gin.Engine {

	// Set gin mode
	if !cfg.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Security middleware
	router.Use(securityMiddleware())

	// CORS middleware
	if cfg.API.EnableCORS {
		router.Use(corsMiddleware())
	}

	// Request logging middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Serve static files
	router.Use(static.Serve("/", static.LocalFile("./public/", true)))

	// Health check endpoint (outside API prefix)
	router.GET("/health", healthController.HealthCheck)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group(cfg.API.Prefix)
	{
		// General API info
		api.GET("", apiController.GetAPIInfo)
		api.GET("/", apiController.GetAPIInfo)

		// Projects routes
		project := api.Group("/project")
		{
			project.GET("", projectController.GetAllProjects)
			project.GET("/stats", projectController.GetProjectsWithStats)
			project.GET("/:id", projectController.GetProjectByID)
			project.POST("", projectController.CreateProject)
			project.PUT("/:id", projectController.UpdateProject)
			project.DELETE("/:id", projectController.DeleteProject)
			project.GET("/:id/services", projectController.GetServicesByProjectID)
		}

		// Services routes
		services := api.Group("/services")
		{
			services.GET("", serviceController.GetAllServices)
		}
	}

	return router
}

// Security middleware
func securityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Expected host validation
		expectedHost := "localhost:8080"
		if c.Request.Host != expectedHost && os.Getenv("SKIP_HOST_CHECK") != "true" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}

		// Security headers
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")

		// Request timestamp
		c.Header("X-Request-Time", time.Now().UTC().Format(time.RFC3339))

		c.Next()
	}
}

// CORS middleware
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
