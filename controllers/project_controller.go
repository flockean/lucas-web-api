package controllers

import (
	"LucasApi/api/models"
	"LucasApi/api/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ProjectController handles HTTP requests for projects
type ProjectController struct {
	projectService *services.ProjectService
}

// NewProjectController creates a new project controller
func NewProjectController(projectService *services.ProjectService) *ProjectController {
	return &ProjectController{
		projectService: projectService,
	}
}

// GetAllProjects handles GET /project
// @Summary Get all projects
// @Description Retrieve all projects from the database
// @Tags project
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]models.Project} "Projects retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to retrieve projects"
// @Router /project [get]
func (c *ProjectController) GetAllProjects(ctx *gin.Context) {
	projects, err := c.projectService.GetAllProjects()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve projects",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Projects retrieved successfully",
		Data:    projects,
	})
}

// GetProjectByID handles GET /project/:id
// @Summary Get project by ID
// @Description Retrieve a specific project by its ID
// @Tags project
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.APIResponse{data=models.Project} "Project retrieved successfully"
// @Failure 400 {object} models.APIResponse "Invalid project ID"
// @Failure 404 {object} models.APIResponse "Project not found"
// @Failure 500 {object} models.APIResponse "Failed to retrieve project"
// @Router /project/{id} [get]
func (c *ProjectController) GetProjectByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid project ID",
			Error:   "ID must be a number",
		})
		return
	}

	project, err := c.projectService.GetProjectByID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "project not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve project",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Project retrieved successfully",
		Data:    project,
	})
}

// CreateProject handles POST /project
// @Summary Create new project
// @Description Create a new project with the provided data
// @Tags project
// @Accept json
// @Produce json
// @Param project body models.CreateProjectRequest true "Project data"
// @Success 201 {object} models.APIResponse{data=models.Project} "Project created successfully"
// @Failure 400 {object} models.APIResponse "Invalid request body or validation error"
// @Failure 500 {object} models.APIResponse "Failed to create project"
// @Router /project [post]
func (c *ProjectController) CreateProject(ctx *gin.Context) {
	var req models.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	project, err := c.projectService.CreateProject(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "project name is required" ||
			err.Error() == "project name must be at least 3 characters long" ||
			err.Error() == "project name must be less than 64 characters" {
			status = http.StatusBadRequest
		}

		ctx.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to create project",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Project created successfully",
		Data:    project,
	})
}

// UpdateProject handles PUT /project/:id
// @Summary Update project
// @Description Update an existing project with the provided data
// @Tags project
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body models.UpdateProjectRequest true "Updated project data"
// @Success 200 {object} models.APIResponse{data=models.Project} "Project updated successfully"
// @Failure 400 {object} models.APIResponse "Invalid project ID or request body"
// @Failure 404 {object} models.APIResponse "Project not found"
// @Failure 500 {object} models.APIResponse "Failed to update project"
// @Router /project/{id} [put]
func (c *ProjectController) UpdateProject(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid project ID",
			Error:   "ID must be a number",
		})
		return
	}

	var req models.UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	project, err := c.projectService.UpdateProject(id, req)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "project not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid project ID" {
			status = http.StatusBadRequest
		}

		ctx.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to update project",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Project updated successfully",
		Data:    project,
	})
}

// DeleteProject handles DELETE /project/:id
// @Summary Delete project
// @Description Delete a project by its ID
// @Tags project
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.APIResponse "Project deleted successfully"
// @Failure 400 {object} models.APIResponse "Invalid project ID"
// @Failure 404 {object} models.APIResponse "Project not found"
// @Failure 409 {object} models.APIResponse "Cannot delete project with active services"
// @Failure 500 {object} models.APIResponse "Failed to delete project"
// @Router /project/{id} [delete]
func (c *ProjectController) DeleteProject(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid project ID",
			Error:   "ID must be a number",
		})
		return
	}

	err = c.projectService.DeleteProject(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "project not found" {
			status = http.StatusNotFound
		} else if err.Error() == "invalid project ID" {
			status = http.StatusBadRequest
		} else if err.Error() == "cannot delete project: active services found" {
			status = http.StatusConflict
		}

		ctx.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to delete project",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Project deleted successfully",
	})
}

// GetProjectsWithStats handles GET /project/stats
// @Summary Get projects with statistics
// @Description Retrieve all projects with additional statistics like service count
// @Tags project
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse{data=[]map[string]interface{}} "Project statistics retrieved successfully"
// @Failure 500 {object} models.APIResponse "Failed to retrieve project statistics"
// @Router /project/stats [get]
func (c *ProjectController) GetProjectsWithStats(ctx *gin.Context) {
	stats, err := c.projectService.GetProjectsWithStats()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve project statistics",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Project statistics retrieved successfully",
		Data:    stats,
	})
}

// GetServicesByProjectID handles GET /project/:id/services
// @Summary Get services by project ID
// @Description Retrieve all services for a specific project
// @Tags project
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} models.APIResponse{data=[]models.Service} "Project services retrieved successfully"
// @Failure 400 {object} models.APIResponse "Invalid project ID"
// @Failure 404 {object} models.APIResponse "Project not found"
// @Failure 500 {object} models.APIResponse "Failed to retrieve project services"
// @Router /project/{id}/services [get]
func (c *ProjectController) GetServicesByProjectID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid project ID",
			Error:   "ID must be a number",
		})
		return
	}

	services, err := c.projectService.GetServicesByProjectID(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err.Error() == "project not found" {
			status = http.StatusNotFound
		}

		ctx.JSON(status, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve project services",
			Error:   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Project services retrieved successfully",
		Data:    services,
	})
}
