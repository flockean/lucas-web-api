package database

import (
	"LucasApi/api/dbutils"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// ProjectRepository handles database operations for projects
type ProjectRepository struct {
	db *DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

// GetAll retrieves all projects from the database
func (r *ProjectRepository) GetAll() ([]dbutils.Project, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at 
		FROM project 
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying projects: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	var projects []dbutils.Project
	for rows.Next() {
		var project dbutils.Project
		err := rows.Scan(&project.ID, &project.Name, &project.Description,
			&project.Status, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning project: %v", err)
			continue
		}
		projects = append(projects, project)
	}

	return projects, nil
}

// GetByID retrieves a project by its ID
func (r *ProjectRepository) GetByID(id int) (*dbutils.Project, error) {
	query := `
		SELECT id, name, description, status, created_at, updated_at 
		FROM project 
		WHERE id = $1
	`
	row := r.db.QueryRow(query, id)

	var project dbutils.Project
	err := row.Scan(&project.ID, &project.Name, &project.Description,
		&project.Status, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("error querying project: %v", err)
	}

	return &project, nil
}

// Create creates a new project
func (r *ProjectRepository) Create(req dbutils.CreateProjectRequest) (*dbutils.Project, error) {
	// Set default status if not provided
	if req.Status == "" {
		req.Status = "active"
	}

	query := `
		INSERT INTO project (name, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, status, created_at, updated_at
	`
	now := time.Now()
	row := r.db.QueryRow(query, req.Name, req.Description, req.Status, now, now)

	var project dbutils.Project
	err := row.Scan(&project.ID, &project.Name, &project.Description,
		&project.Status, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating project: %v", err)
	}

	return &project, nil
}

// Update updates an existing project
func (r *ProjectRepository) Update(id int, req dbutils.UpdateProjectRequest) (*dbutils.Project, error) {
	// Check if project exists
	existingProject, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Build update query dynamically
	updates := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	if req.Name != nil {
		updates = append(updates, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Description != nil {
		updates = append(updates, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Status != nil {
		updates = append(updates, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *req.Status)
		argIndex++
	}

	if len(updates) == 0 {
		return existingProject, nil // No updates requested
	}

	// Add updated_at
	updates = append(updates, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, time.Now())
	argIndex++

	// Add WHERE clause
	args = append(args, id)

	// Build proper query
	var updateFields string
	for i, update := range updates {
		if i > 0 {
			updateFields += ", "
		}
		updateFields += update
	}

	query := fmt.Sprintf(`
		UPDATE project 
		SET %s
		WHERE id = $%d
		RETURNING id, name, description, status, created_at, updated_at
	`, updateFields, argIndex)

	row := r.db.QueryRow(query, args...)

	var project dbutils.Project
	err = row.Scan(&project.ID, &project.Name, &project.Description,
		&project.Status, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("error updating project: %v", err)
	}

	return &project, nil
}

// Delete deletes a project by ID
func (r *ProjectRepository) Delete(id int) error {
	// Check if project exists
	_, err := r.GetByID(id)
	if err != nil {
		return err
	}

	query := `DELETE FROM project WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting project: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

// GetWithStats returns projects with additional statistics
func (r *ProjectRepository) GetWithStats() ([]map[string]interface{}, error) {
	query := `
		SELECT 
			p.id, p.name, p.description, p.status, p.created_at, p.updated_at,
			COUNT(s.id) as service_count
		FROM project p
		LEFT JOIN service s ON p.id = s.project
		GROUP BY p.id, p.name, p.description, p.status, p.created_at, p.updated_at
		ORDER BY p.created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying projects with stats: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	var results []map[string]interface{}
	for rows.Next() {
		var project dbutils.Project
		var serviceCount int

		err := rows.Scan(&project.ID, &project.Name, &project.Description,
			&project.Status, &project.CreatedAt, &project.UpdatedAt, &serviceCount)
		if err != nil {
			log.Printf("Error scanning project with stats: %v", err)
			continue
		}

		projectMap := map[string]interface{}{
			"id":            project.ID,
			"name":          project.Name,
			"description":   project.Description,
			"status":        project.Status,
			"created_at":    project.CreatedAt,
			"updated_at":    project.UpdatedAt,
			"service_count": serviceCount,
		}
		results = append(results, projectMap)
	}

	return results, nil
}

// ServiceRepository handles database operations for services
type ServiceRepository struct {
	db *DB
}

// NewServiceRepository creates a new service repository
func NewServiceRepository(db *DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// GetByProjectID retrieves all services for a specific project
func (r *ServiceRepository) GetByProjectID(projectID int) ([]dbutils.Service, error) {
	query := `
		SELECT id, name, lang, focus, project, created_at, updated_at 
		FROM service 
		WHERE project = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("error querying services: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	var services []dbutils.Service
	for rows.Next() {
		var service dbutils.Service
		err := rows.Scan(&service.ID, &service.Name, &service.Lang,
			&service.Focus, &service.Project, &service.CreatedAt, &service.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning service: %v", err)
			continue
		}
		services = append(services, service)
	}

	return services, nil
}

// GetAll retrieves all services
func (r *ServiceRepository) GetAll() ([]dbutils.Service, error) {
	query := `
		SELECT id, name, lang, focus, project, created_at, updated_at 
		FROM service 
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying services: %v", err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error closing rows: %v", closeErr)
		}
	}()

	var services []dbutils.Service
	for rows.Next() {
		var service dbutils.Service
		err := rows.Scan(&service.ID, &service.Name, &service.Lang,
			&service.Focus, &service.Project, &service.CreatedAt, &service.UpdatedAt)
		if err != nil {
			log.Printf("Error scanning service: %v", err)
			continue
		}
		services = append(services, service)
	}

	return services, nil
}
