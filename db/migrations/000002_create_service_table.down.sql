-- Rollback service table
DROP INDEX IF EXISTS idx_service_created_at;
DROP INDEX IF EXISTS idx_service_lang;
DROP INDEX IF EXISTS idx_service_project;
DROP TABLE IF EXISTS service;
