-- Rollback seed data
DELETE FROM service WHERE name IN ('Authentication Service', 'User Management', 'Dashboard UI');
DELETE FROM project WHERE name IN ('Lucas Web API', 'Frontend Dashboard', 'Mobile App');
