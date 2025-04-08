# Go REST API Server Template

A production-ready REST API server template using Go, Gin, PostgreSQL, and JWT authentication.

## Features

- **REST API** using [Gin framework](https://github.com/gin-gonic/gin)
- **Authentication** with JWT tokens
- **Authorization** middleware with role-based access control
- **PostgreSQL** database with [pgx](https://github.com/jackc/pgx) driver
- **Raw SQL** queries (no ORM)
- **Structured logging** with [zerolog](https://github.com/rs/zerolog)
- **Configuration** using [Viper](https://github.com/spf13/viper)
- **Docker** support with multi-stage builds and distroless images
- **SOLID** principles and clean architecture

## Project Structure

```
my-app/
├── cmd/
│   └── server/            # Application entry point
│       └── main.go
├── internal/
│   ├── api/               # API layer
│   │   ├── handlers/      # Request handlers
│   │   ├── middleware/    # HTTP middleware
│   │   └── routes/        # Route definitions
│   ├── config/            # Configuration
│   ├── domain/            # Domain layer
│   │   ├── models/        # Data structures
│   │   └── services/      # Business logic
│   ├── db/                # Database layer
│   │   ├── postgres/      # Postgres connection
│   │   └── migrations/    # SQL migration files
│   └── utils/             # Utility functions
├── Dockerfile             # Docker image definition
├── docker-compose.yml     # Docker services configuration
├── .dockerignore          # Docker build exclusions
├── .gitignore             # Git exclusions
├── go.mod                 # Go module definition
└── config.yaml            # Application configuration
```

## Getting Started

### Prerequisites

- Go 1.23+
- PostgreSQL
- Docker & Docker Compose (optional)

### Local Development Setup

1. Clone the repository:

   ```bash
   git clone https://github.com/ramhdi/go-backend-starter.git
   cd go-backend-starter
   ```

2. Install dependencies:

   ```bash
   go mod download
   ```

3. Set up the database:

   ```bash
   # Using Docker
   docker-compose up -d postgres

   # Or connect to your existing PostgreSQL instance
   ```

4. Run the database migrations:

   ```bash
   # Apply migrations manually using the SQL files in internal/db/migrations
   psql -U postgres -d myapp -f internal/db/migrations/001_create_users_table.sql
   ```

5. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

### Using Docker

1. Build the Docker image:

   ```bash
   docker build -t go-backend-starter:latest .
   ```

2. Run with Docker Compose:
   ```bash
   # Uncomment the app service in docker-compose.yml first
   docker-compose up -d
   ```

## API Endpoints

### Authentication

- `POST /api/auth/login` - Login with username and password

### Users (Admin only)

- `GET /api/users` - List all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create a new user
- `PUT /api/users/:id` - Update a user
- `DELETE /api/users/:id` - Delete a user

### Current User

- `GET /api/me` - Get current user information

## Configuration

The application can be configured using:

1. Configuration file (`config.yaml`)
2. Environment variables

### Environment Variables

| Variable           | Description                          | Default              |
| ------------------ | ------------------------------------ | -------------------- |
| SERVER_PORT        | HTTP server port                     | 8080                 |
| SERVER_ENVIRONMENT | Environment (development/production) | development          |
| DATABASE_HOST      | PostgreSQL host                      | localhost            |
| DATABASE_PORT      | PostgreSQL port                      | 5432                 |
| DATABASE_USER      | PostgreSQL username                  | postgres             |
| DATABASE_PASSWORD  | PostgreSQL password                  | postgres             |
| DATABASE_DBNAME    | PostgreSQL database name             | myapp                |
| DATABASE_SSLMODE   | PostgreSQL SSL mode                  | disable              |
| JWT_SECRET         | Secret key for JWT signing           | your-secret-key-here |
| JWT_EXPIRATION     | JWT token expiration (minutes)       | 60                   |

## Default Admin User

The migration script creates a default admin user:

- **Username**: admin
- **Password**: admin
- **Email**: admin@example.com
- **Role**: admin

> **Note**: Change these credentials in production!

## License

This project is licensed under the MIT License - see the LICENSE file for details.
