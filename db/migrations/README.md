# Database Migrations

This directory contains all database migrations for the Lucas Web API.

## Migration File Naming Convention

Migrations follow this naming pattern:
```
{version}_{description}.{direction}.sql
```

Example:
- `000001_create_project_table.up.sql` - Creates the project table
- `000001_create_project_table.down.sql` - Rolls back the project table

## Creating a New Migration

1. Determine the next version number (look at the highest existing version and increment)
2. Create two files:
   - `{version}_{description}.up.sql` - Contains the migration SQL
   - `{version}_{description}.down.sql` - Contains the rollback SQL

Example:
```bash
# Create migration files for adding a new column
touch db/migrations/000005_add_email_to_project.up.sql
touch db/migrations/000005_add_email_to_project.down.sql
```

### Up Migration Example (`000005_add_email_to_project.up.sql`):
```sql
ALTER TABLE project ADD COLUMN email VARCHAR(255);
CREATE INDEX IF NOT EXISTS idx_project_email ON project(email);
```

### Down Migration Example (`000005_add_email_to_project.down.sql`):
```sql
DROP INDEX IF EXISTS idx_project_email;
ALTER TABLE project DROP COLUMN IF EXISTS email;
```

## How Migrations Work

1. **Automatic Execution**: Migrations run automatically when the application starts
2. **Tracking**: Applied migrations are tracked in the `schema_migrations` table
3. **Idempotent**: Migrations can be run multiple times safely (already applied migrations are skipped)
4. **Order**: Migrations are applied in version order (000001, 000002, 000003, etc.)
5. **Dirty Flag**: If a migration fails, it's marked as "dirty" to prevent further migrations until fixed

## Current Migrations

- **000001**: Creates the schema_migrations tracking table
- **000002**: Creates the project table with indexes
- **000003**: Creates the service table with foreign key to project
- **000004**: Seeds initial sample data

## Migration States

Each migration in the `schema_migrations` table has:
- `version`: The migration version number
- `dirty`: Boolean indicating if the migration failed (true = failed, false = success)
- `applied_at`: Timestamp when the migration was successfully applied

## Best Practices

1. **Never edit applied migrations**: Always create a new migration to modify the schema
2. **Test rollbacks**: Always test the down migration to ensure it works
3. **Keep migrations small**: Break large changes into multiple migrations
4. **Use transactions**: Wrap complex migrations in transactions when possible
5. **Document changes**: Use clear, descriptive names for migrations
6. **Avoid data migrations**: Separate schema changes from data migrations when possible

## Troubleshooting

### Migration marked as dirty
If a migration fails and is marked as dirty:
1. Fix the migration SQL
2. Manually update the schema_migrations table:
   ```sql
   DELETE FROM schema_migrations WHERE version = {failed_version};
   ```
3. Restart the application to retry the migration

### Reset database (development only)
To start fresh:
```bash
docker compose down -v
docker compose up -d db
# Restart your application
```
