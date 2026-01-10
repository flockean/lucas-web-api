package controllers

import (
	"LucasApi/api/models"
	"LucasApi/api/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ServiceController handles HTTP requests for services
type ServiceController struct {
	serviceService *services.ServiceService
}

// NewServiceController creates a new service controller
func NewServiceController(serviceService *services.ServiceService) *ServiceController {
	return &ServiceController{
		serviceService: serviceService,
	}
}

// GetAllServices handles GET /services
// @Summary Get all services
// @Description Retrieve all services from the database
// @Tags services
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.Service} "Services retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to retrieve services"
// @Router /services [get]
func (c *ServiceController) GetAllServices(ctx *gin.Context) {
	services, err := c.serviceService.GetAllServices()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve services",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Services retrieved successfully",
		Data:    services,
	})
}

// HealthController handles health check requests
type HealthController struct{}

// NewHealthController creates a new health controller
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck handles GET /health
// @Summary Health check
// @Description Check if the service is healthy and running
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=map[string]interface{}} "Service is healthy"
// @Router /health [get]
func (c *HealthController) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Service is healthy",
		Data: map[string]interface{}{
			"status":    "healthy",
			"version":   "1.0.0",
			"timestamp": ctx.GetString("X-Request-Time"),
		},
	})
}

// APIController handles general API requests
type APIController struct{}

// NewAPIController creates a new API controller
func NewAPIController() *APIController {
	return &APIController{}
}

// GetAPIInfo handles GET /api
// @Summary API Information
// @Description Get general information about the API
// @Tags system
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=map[string]interface{}} "API information retrieved successfully"
// @Router / [get]
func (c *APIController) GetAPIInfo(ctx *gin.Context) {
	// Check if user is authenticated
	user, authenticated := ctx.Get("user")

	data := map[string]interface{}{
		"name":        "Lucas Web API",
		"version":     "1.0.0",
		"description": "REST API for project and service management with OAuth2",
		"endpoints": map[string]interface{}{
			"project":  "/api/project",
			"services": "/api/services",
			"health":   "/health",
			"auth": map[string]interface{}{
				"login":    "/api/auth/login",
				"logout":   "/api/auth/logout",
				"status":   "/api/auth/status",
				"me":       "/api/auth/me",
				"callback": "/api/auth/callback",
			},
		},
		"oauth2": map[string]interface{}{
			"enabled":  ctx.GetBool("oauth2_enabled"),
			"provider": ctx.GetString("oauth2_provider"),
		},
		"authenticated": authenticated,
	}

	if authenticated {
		data["user"] = user
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Lucas Web API - OAuth2 Ready",
		Data:    data,
	})
}
