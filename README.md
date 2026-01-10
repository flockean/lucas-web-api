# Lucas Web API

A REST API built with Go and Gin framework featuring OAuth2 authentication and PostgreSQL database.

## Features

- RESTful API with Gin framework
- OAuth2 authentication (GitHub, Google, Azure)
- PostgreSQL database integration
- Session-based authentication
- Swagger documentation
- Docker containerization
- CORS support

## Quick Start

### Prerequisites

- Go 1.24+
- PostgreSQL
- Docker (optional)

### Installation

1. Clone the repository
```bash
git clone <repository-url>
cd lucas-web-api
```

2. Copy environment file
```bash
cp .env.example .env
```

3. Configure OAuth2 credentials in `.env`
```bash
OAUTH2_ENABLED=true
OAUTH2_PROVIDER=github
OAUTH2_CLIENT_ID=your_client_id
OAUTH2_CLIENT_SECRET=your_client_secret
```

4. Start with Docker
```bash
docker-compose up -d
```

Or run locally:
```bash
go run .
```

## API Endpoints

### Authentication
- `GET /api/auth/login` - OAuth2 login
- `GET /api/auth/callback` - OAuth2 callback
- `POST /api/auth/logout` - Logout
- `GET /api/auth/login-url` - Get login URL

### Projects
- `GET /api/project` - List all projects
- `POST /api/project` - Create project (authenticated)
- `GET /api/project/{id}` - Get project by ID
- `PUT /api/project/{id}` - Update project (authenticated)
- `DELETE /api/project/{id}` - Delete project (authenticated)

### Services
- `GET /api/service` - List all services
- `GET /api/service/project/{id}` - List services by project

### Documentation
- `GET /swagger/index.html` - Swagger UI
- `GET /docs/index.html` - Alternative docs endpoint

## Configuration

Environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `OAUTH2_ENABLED` | Enable OAuth2 authentication | `false` |
| `OAUTH2_PROVIDER` | OAuth2 provider (github/google/azure) | `github` |
| `OAUTH2_CLIENT_ID` | OAuth2 client ID | - |
| `OAUTH2_CLIENT_SECRET` | OAuth2 client secret | - |
| `OAUTH2_REDIRECT_URL` | OAuth2 redirect URL | `http://localhost:8080/api/auth/callback` |
| `API_ENABLE_AUTH` | Require auth for write operations | `true` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `lucas_api` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | - |

## OAuth2 Setup

### GitHub
1. Go to GitHub Settings > Developer settings > OAuth Apps
2. Create new OAuth App
3. Set Authorization callback URL: `http://localhost:8080/api/auth/callback`
4. Copy Client ID and Client Secret to `.env`

### Google
1. Go to Google Cloud Console
2. Create OAuth 2.0 credentials
3. Add authorized redirect URI: `http://localhost:8080/api/auth/callback`
4. Copy credentials to `.env`

## Development

### Generate Swagger docs

Note: Might not without config, because because of binary from go
```bash
swag init
```

### Run tests
```bash
go test ./...
```

### Lint code
```bash
golangci-lint run
```

## Docker

### Build image
```bash
docker build -t lucas-web-api .
```

### Run with docker-compose
```bash
docker-compose up -d
```

## Project Structure

```
.
├── clients/          # API clients
├── config/           # Configuration management
├── controllers/      # HTTP handlers
├── database/         # Database layer
├── db/              # Database migrations
├── docs/            # Swagger documentation
├── middleware/      # HTTP middleware
├── models/          # Data models
├── services/        # Business logic
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── main.go
```

## License

MIT