-- Rollback project table
DROP INDEX IF EXISTS idx_project_created_at;
DROP INDEX IF EXISTS idx_project_status;
DROP TABLE IF EXISTS project CASCADE;
