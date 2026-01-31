package controllers

import (
	"LucasApi/api/models"
	"LucasApi/api/services"
	"net/http"
	"strconv"

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

// GetServiceByID handles GET /services/:id
// @Summary Get service by ID
// @Description Retrieve a specific service by its ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} models.APIResponse{data=models.Service} "Service retrieved successfully"
// @Failure 400 {object} models.APIResponse "Invalid service ID"
// @Failure 404 {object} models.APIResponse "Service not found"
// @Failure 500 {object} models.APIResponse "Failed to retrieve service"
// @Router /services/{id} [get]
func (c *ServiceController) GetServiceByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid service ID",
			Error:   "ID must be a number",
		})
		return
	}

	service, err := c.serviceService.GetServiceByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "service not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve service",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Service retrieved successfully",
		Data:    service,
	})
}

// CreateService handles POST /services
// @Summary Create new service
// @Description Create a new service with the provided data
// @Tags services
// @Accept json
// @Produce json
// @Param service body models.CreateServiceRequest true "Service data"
// @Success 201 {object} models.APIResponse{data=models.Service} "Service created successfully"
// @Failure 400 {object} models.APIResponse "Invalid request body or validation error"
// @Failure 500 {object} models.APIResponse "Failed to create service"
// @Router /services [post]
func (c *ServiceController) CreateService(ctx *gin.Context) {
	var req models.CreateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	service, err := c.serviceService.CreateService(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "service name is required" ||
			err.Error() == "valid project ID is required" {
			status = http.StatusBadRequest
		}

		ctx.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to create service",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Service created successfully",
		Data:    service,
	})
}

// UpdateService handles PUT /services/:id
// @Summary Update service
// @Description Update an existing service with the provided data
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Param service body models.UpdateServiceRequest true "Updated service data"
// @Success 200 {object} models.APIResponse{data=models.Service} "Service updated successfully"
// @Failure 400 {object} models.APIResponse "Invalid service ID or request body"
// @Failure 404 {object} models.APIResponse "Service not found"
// @Failure 500 {object} models.APIResponse "Failed to update service"
// @Router /services/{id} [put]
func (c *ServiceController) UpdateService(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid service ID",
			Error:   "ID must be a number",
		})
		return
	}

	var req models.UpdateServiceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	service, err := c.serviceService.UpdateService(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "service not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update service",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Service updated successfully",
		Data:    service,
	})
}

// DeleteService handles DELETE /services/:id
// @Summary Delete service
// @Description Delete a service by its ID
// @Tags services
// @Accept json
// @Produce json
// @Param id path int true "Service ID"
// @Success 200 {object} models.APIResponse "Service deleted successfully"
// @Failure 400 {object} models.APIResponse "Invalid service ID"
// @Failure 404 {object} models.APIResponse "Service not found"
// @Failure 500 {object} models.APIResponse "Failed to delete service"
// @Router /services/{id} [delete]
func (c *ServiceController) DeleteService(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid service ID",
			Error:   "ID must be a number",
		})
		return
	}

	err = c.serviceService.DeleteService(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "service not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to delete service",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Service deleted successfully",
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
