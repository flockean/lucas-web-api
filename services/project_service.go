package services

import (
	"LucasApi/api/database"
	"LucasApi/api/models"
	"fmt"
	"log"
)

// ProjectService contains business logic for projects
type ProjectService struct {
	projectRepo *database.ProjectRepository
	serviceRepo *database.ServiceRepository
}

// NewProjectService creates a new project service
func NewProjectService(projectRepo *database.ProjectRepository, serviceRepo *database.ServiceRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		serviceRepo: serviceRepo,
	}
}

// GetAllProjects retrieves all projects
func (s *ProjectService) GetAllProjects() ([]models.Project, error) {
	log.Println("Service: Getting all projects")

	projects, err := s.projectRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %v", err)
	}

	log.Printf("Service: Retrieved %d projects", len(projects))
	return projects, nil
}

// GetProjectByID retrieves a project by ID
func (s *ProjectService) GetProjectByID(id int) (*models.Project, error) {
	log.Printf("Service: Getting project with ID %d", id)

	if id <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", id)
	}

	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get project: %v", err)
	}

	log.Printf("Service: Retrieved project '%s'", project.Name)
	return project, nil
}

// CreateProject creates a new project
func (s *ProjectService) CreateProject(req models.CreateProjectRequest) (*models.Project, error) {
	log.Printf("Service: Creating project '%s'", req.Name)

	// Business logic validation
	if err := s.validateCreateProjectRequest(req); err != nil {
		return nil, err
	}

	project, err := s.projectRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %v", err)
	}

	log.Printf("Service: Created project '%s' with ID %d", project.Name, project.ID)
	return project, nil
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(id int, req models.UpdateProjectRequest) (*models.Project, error) {
	log.Printf("Service: Updating project with ID %d", id)

	if id <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", id)
	}

	// Business logic validation
	if err := s.validateUpdateProjectRequest(req); err != nil {
		return nil, err
	}

	project, err := s.projectRepo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update project: %v", err)
	}

	log.Printf("Service: Updated project '%s'", project.Name)
	return project, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(id int) error {
	log.Printf("Service: Deleting project with ID %d", id)

	if id <= 0 {
		return fmt.Errorf("invalid project ID: %d", id)
	}

	// Business logic: Check if project has active services
	services, err := s.serviceRepo.GetByProjectID(id)
	if err != nil {
		return fmt.Errorf("failed to check project services: %v", err)
	}

	if len(services) > 0 {
		return fmt.Errorf("cannot delete project: %d active services found", len(services))
	}

	err = s.projectRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %v", err)
	}

	log.Printf("Service: Deleted project with ID %d", id)
	return nil
}

// GetProjectsWithStats retrieves projects with statistics
func (s *ProjectService) GetProjectsWithStats() ([]map[string]interface{}, error) {
	log.Println("Service: Getting projects with statistics")

	stats, err := s.projectRepo.GetWithStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get project stats: %v", err)
	}

	log.Printf("Service: Retrieved statistics for %d projects", len(stats))
	return stats, nil
}

// GetServicesByProjectID retrieves services for a project
func (s *ProjectService) GetServicesByProjectID(projectID int) ([]models.Service, error) {
	log.Printf("Service: Getting services for project ID %d", projectID)

	if projectID <= 0 {
		return nil, fmt.Errorf("invalid project ID: %d", projectID)
	}

	// Ensure project exists
	_, err := s.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, fmt.Errorf("project not found: %v", err)
	}

	services, err := s.serviceRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	log.Printf("Service: Retrieved %d services for project %d", len(services), projectID)
	return services, nil
}

// Business logic validation methods

func (s *ProjectService) validateCreateProjectRequest(req models.CreateProjectRequest) error {
	if req.Name == "" {
		return fmt.Errorf("project name is required")
	}

	if len(req.Name) < 3 {
		return fmt.Errorf("project name must be at least 3 characters long")
	}

	if len(req.Name) > 64 {
		return fmt.Errorf("project name must be less than 64 characters")
	}

	if req.Status != "" {
		validStatuses := map[string]bool{
			"active":    true,
			"inactive":  true,
			"planned":   true,
			"completed": true,
			"archived":  true,
		}
		if !validStatuses[req.Status] {
			return fmt.Errorf("invalid project status: %s", req.Status)
		}
	}

	return nil
}

func (s *ProjectService) validateUpdateProjectRequest(req models.UpdateProjectRequest) error {
	if req.Name != nil {
		if *req.Name == "" {
			return fmt.Errorf("project name cannot be empty")
		}
		if len(*req.Name) < 3 {
			return fmt.Errorf("project name must be at least 3 characters long")
		}
		if len(*req.Name) > 64 {
			return fmt.Errorf("project name must be less than 64 characters")
		}
	}

	if req.Status != nil {
		validStatuses := map[string]bool{
			"active":    true,
			"inactive":  true,
			"planned":   true,
			"completed": true,
			"archived":  true,
		}
		if !validStatuses[*req.Status] {
			return fmt.Errorf("invalid project status: %s", *req.Status)
		}
	}

	return nil
}

// ServiceService contains business logic for services
type ServiceService struct {
	serviceRepo *database.ServiceRepository
	projectRepo *database.ProjectRepository
}

// NewServiceService creates a new service service
func NewServiceService(serviceRepo *database.ServiceRepository, projectRepo *database.ProjectRepository) *ServiceService {
	return &ServiceService{
		serviceRepo: serviceRepo,
		projectRepo: projectRepo,
	}
}

// GetAllServices retrieves all services
func (s *ServiceService) GetAllServices() ([]models.Service, error) {
	log.Println("Service: Getting all services")

	services, err := s.serviceRepo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get services: %v", err)
	}

	log.Printf("Service: Retrieved %d services", len(services))
	return services, nil
}

// GetServiceByID retrieves a service by ID
func (s *ServiceService) GetServiceByID(id int) (*models.Service, error) {
	log.Printf("Service: Getting service with ID %d", id)

	if id <= 0 {
		return nil, fmt.Errorf("invalid service ID: %d", id)
	}

	service, err := s.serviceRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %v", err)
	}

	log.Printf("Service: Retrieved service '%s'", service.Name)
	return service, nil
}

// CreateService creates a new service
func (s *ServiceService) CreateService(req models.CreateServiceRequest) (*models.Service, error) {
	log.Printf("Service: Creating service '%s'", req.Name)

	// Business logic validation
	if err := s.validateCreateServiceRequest(req); err != nil {
		return nil, err
	}

	// Verify that the project exists
	_, err := s.projectRepo.GetByID(req.Project)
	if err != nil {
		return nil, fmt.Errorf("invalid project: %v", err)
	}

	service, err := s.serviceRepo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %v", err)
	}

	log.Printf("Service: Created service '%s' with ID %d", service.Name, service.ID)
	return service, nil
}

// UpdateService updates an existing service
func (s *ServiceService) UpdateService(id int, req models.UpdateServiceRequest) (*models.Service, error) {
	log.Printf("Service: Updating service with ID %d", id)

	if id <= 0 {
		return nil, fmt.Errorf("invalid service ID: %d", id)
	}

	// Business logic validation
	if err := s.validateUpdateServiceRequest(req); err != nil {
		return nil, err
	}

	// If project is being updated, verify it exists
	if req.Project != nil {
		_, err := s.projectRepo.GetByID(*req.Project)
		if err != nil {
			return nil, fmt.Errorf("invalid project: %v", err)
		}
	}

	service, err := s.serviceRepo.Update(id, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %v", err)
	}

	log.Printf("Service: Updated service '%s'", service.Name)
	return service, nil
}

// DeleteService deletes a service
func (s *ServiceService) DeleteService(id int) error {
	log.Printf("Service: Deleting service with ID %d", id)

	if id <= 0 {
		return fmt.Errorf("invalid service ID: %d", id)
	}

	err := s.serviceRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete service: %v", err)
	}

	log.Printf("Service: Deleted service with ID %d", id)
	return nil
}

// Business logic validation methods for services

func (s *ServiceService) validateCreateServiceRequest(req models.CreateServiceRequest) error {
	if req.Name == "" {
		return fmt.Errorf("service name is required")
	}

	if len(req.Name) < 2 {
		return fmt.Errorf("service name must be at least 2 characters long")
	}

	if len(req.Name) > 64 {
		return fmt.Errorf("service name must be less than 64 characters")
	}

	if req.Project <= 0 {
		return fmt.Errorf("valid project ID is required")
	}

	return nil
}

func (s *ServiceService) validateUpdateServiceRequest(req models.UpdateServiceRequest) error {
	if req.Name != nil {
		if *req.Name == "" {
			return fmt.Errorf("service name cannot be empty")
		}
		if len(*req.Name) < 2 {
			return fmt.Errorf("service name must be at least 2 characters long")
		}
		if len(*req.Name) > 64 {
			return fmt.Errorf("service name must be less than 64 characters")
		}
	}

	if req.Project != nil && *req.Project <= 0 {
		return fmt.Errorf("project ID must be valid")
	}

	return nil
}
