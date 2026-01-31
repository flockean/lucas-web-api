-- Create service table
CREATE TABLE IF NOT EXISTS service (
    id SERIAL PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    lang VARCHAR(64),
    focus VARCHAR(64),
    project INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_service_project FOREIGN KEY (project) 
        REFERENCES project(id) 
        ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_service_project ON service(project);
CREATE INDEX IF NOT EXISTS idx_service_lang ON service(lang);
CREATE INDEX IF NOT EXISTS idx_service_created_at ON service(created_at DESC);
