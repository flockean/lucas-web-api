package clients

import (
	"LucasApi/api/config"
	"LucasApi/api/dbutils"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// ProjectsClient is a REST client for Projects API
type ProjectsClient struct {
	client *resty.Client
	config *config.ClientConfig
}

// NewProjectsClient creates a new Projects REST client with config
func NewProjectsClient(cfg *config.ClientConfig) *ProjectsClient {
	client := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(cfg.Timeout).
		SetHeader("User-Agent", cfg.UserAgent).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	if cfg.Debug {
		client.SetDebug(true)
	}

	// Set retry mechanism
	client.SetRetryCount(cfg.RetryCount).
		SetRetryWaitTime(cfg.RetryWaitTime).
		SetRetryMaxWaitTime(cfg.RetryWaitTime * 10)

	return &ProjectsClient{
		client: client,
		config: cfg,
	}
}

// GetAllProjects retrieves all projects
func (pc *ProjectsClient) GetAllProjects() ([]dbutils.Project, error) {
	var response dbutils.APIResponse

	resp, err := pc.client.R().
		SetResult(&response).
		SetError(&response).
		Get("/projects")

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode(), response.Error)
	}

	if !response.Success {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	// Convert response.Data to []dbutils.Project
	projects, ok := response.Data.([]dbutils.Project)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return projects, nil
}

// GetProjectByID retrieves a project by ID
func (pc *ProjectsClient) GetProjectByID(id int) (*dbutils.Project, error) {
	var response dbutils.APIResponse

	resp, err := pc.client.R().
		SetResult(&response).
		SetError(&response).
		Get(fmt.Sprintf("/projects/%d", id))

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode(), response.Error)
	}

	if !response.Success {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	// Convert response.Data to *dbutils.Project
	project, ok := response.Data.(*dbutils.Project)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return project, nil
}

// CreateProject creates a new project
func (pc *ProjectsClient) CreateProject(req dbutils.CreateProjectRequest) (*dbutils.Project, error) {
	var response dbutils.APIResponse

	resp, err := pc.client.R().
		SetBody(req).
		SetResult(&response).
		SetError(&response).
		Post("/projects")

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode(), response.Error)
	}

	if !response.Success {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	// Convert response.Data to *dbutils.Project
	project, ok := response.Data.(*dbutils.Project)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return project, nil
}

// UpdateProject updates an existing project
func (pc *ProjectsClient) UpdateProject(id int, req dbutils.UpdateProjectRequest) (*dbutils.Project, error) {
	var response dbutils.APIResponse

	resp, err := pc.client.R().
		SetBody(req).
		SetResult(&response).
		SetError(&response).
		Put(fmt.Sprintf("/projects/%d", id))

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode(), response.Error)
	}

	if !response.Success {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	// Convert response.Data to *dbutils.Project
	project, ok := response.Data.(*dbutils.Project)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return project, nil
}

// DeleteProject deletes a project by ID
func (pc *ProjectsClient) DeleteProject(id int) error {
	var response dbutils.APIResponse

	resp, err := pc.client.R().
		SetResult(&response).
		SetError(&response).
		Delete(fmt.Sprintf("/projects/%d", id))

	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("API error (%d): %s", resp.StatusCode(), response.Error)
	}

	if !response.Success {
		return fmt.Errorf("API returned error: %s", response.Message)
	}

	return nil
}

// GetProjectsWithStats retrieves projects with statistics
func (pc *ProjectsClient) GetProjectsWithStats() ([]map[string]interface{}, error) {
	var response dbutils.APIResponse

	resp, err := pc.client.R().
		SetResult(&response).
		SetError(&response).
		Get("/projects/stats")

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode(), response.Error)
	}

	if !response.Success {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	// Convert response.Data to []map[string]interface{}
	stats, ok := response.Data.([]map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return stats, nil
}

// GetServicesByProjectID retrieves services for a specific project
func (pc *ProjectsClient) GetServicesByProjectID(projectID int) ([]dbutils.Service, error) {
	var response dbutils.APIResponse

	resp, err := pc.client.R().
		SetResult(&response).
		SetError(&response).
		Get(fmt.Sprintf("/projects/%d/services", projectID))

	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode(), response.Error)
	}

	if !response.Success {
		return nil, fmt.Errorf("API returned error: %s", response.Message)
	}

	// Convert response.Data to []dbutils.Service
	services, ok := response.Data.([]dbutils.Service)
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return services, nil
}
