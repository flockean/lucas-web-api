CREATE TABLE IF NOT EXISTS project(
    id      SERIAL PRIMARY KEY,
    name            VARCHAR(64) NOT NULL,
    description     TEXT,
    status          VARCHAR(32) DEFAULT 'active',
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS service(
    id      SERIAL PRIMARY KEY,
    name            VARCHAR(64) NOT NULL,
    lang            VARCHAR(64),
    focus           VARCHAR(64),
    project         INT,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project) REFERENCES project(id)
    ON DELETE CASCADE
);

-- Insert sample data #TODO: Remove someday
INSERT INTO project (name, description, status) VALUES 
('Lucas Web API', 'Main API project for web services', 'active'),
('Frontend Dashboard', 'React-based dashboard for project management', 'active'),
('Mobile App', 'Mobile companion app', 'planned')
ON CONFLICT DO NOTHING;

INSERT INTO service (name, lang, focus, project) VALUES 
('Authentication Service', 'Go', 'Security', 1),
('User Management', 'Go', 'Backend', 1),
('Dashboard UI', 'React', 'Frontend', 2)
ON CONFLICT DO NOTHING;
