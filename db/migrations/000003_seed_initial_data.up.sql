-- Insert sample data for projects
INSERT INTO project (name, description, status) VALUES 
('Lucas Web API', 'Main API project for web services', 'active'),
('Frontend Dashboard', 'React-based dashboard for project management', 'active'),
('Mobile App', 'Mobile companion app', 'planned')
ON CONFLICT DO NOTHING;

-- Insert sample data for services (only if projects exist)
INSERT INTO service (name, lang, focus, project) 
SELECT 'Authentication Service', 'Go', 'Security', p.id 
FROM project p WHERE p.name = 'Lucas Web API'
ON CONFLICT DO NOTHING;

INSERT INTO service (name, lang, focus, project) 
SELECT 'User Management', 'Go', 'Backend', p.id 
FROM project p WHERE p.name = 'Lucas Web API'
ON CONFLICT DO NOTHING;

INSERT INTO service (name, lang, focus, project) 
SELECT 'Dashboard UI', 'React', 'Frontend', p.id 
FROM project p WHERE p.name = 'Frontend Dashboard'
ON CONFLICT DO NOTHING;
