# Lucas Web API

A REST API for project and service management built with Go, featuring OAuth2 authentication and PostgreSQL database.

## Features

- RESTful API with full CRUD operations
- OAuth2 authentication (GitHub/Google)
- PostgreSQL with automatic migrations
- Session-based authentication
- Swagger API documentation
- Public read endpoints, authenticated write operations

## Quick Start

1. **Clone and configure**
```bash
git clone <repository-url>
cd lucas-web-api
cp .env.example .env
```

2. **Configure environment variables** (see Configuration section below)

3. **Start the application**
```bash
docker compose up -d
```

The API will be available at `http://localhost:8080`

**API Documentation:** http://localhost:8080/swagger/index.html

## API Endpoints

### Public Endpoints (No Authentication Required)
- `GET /api/project` - List all projects
- `GET /api/project/stats` - Get project statistics
- `GET /api/project/:id` - Get project by ID
- `GET /api/project/:id/services` - Get services for a project
- `GET /api/services` - List all services
- `GET /api/services/:id` - Get service by ID

### Authenticated Endpoints (OAuth2 Required)
- `POST /api/project` - Create project
- `PUT /api/project/:id` - Update project
- `DELETE /api/project/:id` - Delete project
- `POST /api/services` - Create service
- `PUT /api/services/:id` - Update service
- `DELETE /api/services/:id` - Delete service

### Authentication Endpoints
- `GET /api/auth/login` - Start OAuth2 login flow
- `GET /api/auth/callback` - OAuth2 callback handler
- `POST /api/auth/logout` - Logout user
- `GET /api/auth/status` - Check authentication status
- `GET /api/auth/me` - Get current user info

## Configuration

Create a `.env` file with these variables:

```bash
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_DEBUG=true

# Database
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=rocket123
DB_NAME=projectDb

# OAuth2 (optional - enables authentication)
OAUTH2_ENABLED=true
OAUTH2_PROVIDER=github
OAUTH2_CLIENT_ID=your_client_id_here
OAUTH2_CLIENT_SECRET=your_client_secret_here
OAUTH2_REDIRECT_URL=http://localhost:8080/api/auth/callback
OAUTH2_SESSION_SECRET=your_random_secret_key

# API Settings
API_ENABLE_AUTH=true
API_ENABLE_CORS=true
```

## OAuth2 Setup

### GitHub OAuth App
1. Go to [GitHub Settings > Developer settings > OAuth Apps](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in:
   - **Application name:** Lucas Web API
   - **Homepage URL:** `http://localhost:8080`
   - **Authorization callback URL:** `http://localhost:8080/api/auth/callback`
4. Copy the **Client ID** and generate a **Client Secret**
5. Add them to your `.env` file

### Google OAuth
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Navigate to APIs & Services > Credentials
4. Create OAuth 2.0 Client ID
5. Add authorized redirect URI: `http://localhost:8080/api/auth/callback`
6. Copy credentials to `.env` and set `OAUTH2_PROVIDER=google`

## Database Migrations

Migrations run automatically on startup. Migration files are located in `db/migrations/`.

To create a new migration:
```bash
# Create migration files (replace 000004 with next version number)
touch db/migrations/000004_description.up.sql
touch db/migrations/000004_description.down.sql
```

## Running Locally Without Docker

```bash
# Start PostgreSQL
docker compose up -d db

# Install dependencies
go mod download

# Run the application
go run main.go
```

## License

MIT
