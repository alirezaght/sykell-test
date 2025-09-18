# Sykell Backend

A Go backend application built with Echo framework, SQLC for type-safe SQL queries, and golang-migrate for database migrations with MySQL.

## Features

- **Echo Framework**: Fast and minimalist web framework
- **SQLC**: Generate type-safe Go code from SQL
- **Database Migrations**: Version control for database schema using golang-migrate
- **MySQL**: Primary database with Docker support
- **Docker Compose**: Easy development environment setup
- **Environment Configuration**: Flexible configuration management
- **RESTful API**: Well-structured REST endpoints
- **PhpMyAdmin**: Web-based MySQL administration tool

## Project Structure

```
backend/
├── cmd/                    # Application entrypoints
│   └── main.go            # Main application
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── feature1/             # This is a feature first SOLID friendly layout
│   ├── feature2/             
│   └── db/              # Generated SQLC code (after generation)
├── migrations/           # Database migration files
├── sql/                 # SQL files for SQLC
│   ├── queries/         # SQL queries
│   └── schema/          # Database schema
├── init-scripts/        # MySQL initialization scripts
├── docker-compose.yml   # Docker services configuration
├── .env.example         # Environment variables example
├── sqlc.yaml           # SQLC configuration
└── Makefile           # Build and development commands
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional, for using Makefile commands)

### Quick Start with Docker

1. Clone the repository and navigate to the backend directory
2. Start the complete development environment:
   ```bash
   make dev-setup
   ```
   This will:
   - Start MySQL and PhpMyAdmin containers
   - Copy environment configuration
   - Install development tools
   - Run database migrations
   - Generate SQLC code

3. Start the application:
   ```bash
   make run
   ```

### Manual Installation

1. Start the database:
   ```bash
   make docker-up
   ```

2. Install development tools and setup environment:
   ```bash
   make setup
   ```

3. Update the `.env` file with your database configuration
4. Install dependencies:
   ```bash
   go mod tidy
   ```

5. Generate SQLC code:
   ```bash
   make sqlc-generate
   ```

6. Run database migrations:
   ```bash
   make migrate-up
   ```

7. Start the development server:
   ```bash
   make dev
   # or
   make run
   ```

## API Endpoints

### Health Check
- `GET /` - API information
- `GET /health` - Health check

### Users
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create new user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

## Development

### Available Make Commands

#### Application Commands
- `make build` - Build the application
- `make run` - Build and run the application
- `make dev` - Run with hot reload (requires air)
- `make test` - Run tests
- `make clean` - Clean build artifacts
- `make fmt` - Format code
- `make tidy` - Tidy dependencies

#### Docker Commands
- `make docker-up` - Start MySQL and PhpMyAdmin containers
- `make docker-down` - Stop all containers
- `make docker-logs` - View container logs
- `make docker-clean` - Stop containers and remove volumes

#### Database Commands
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback database migrations
- `make migrate-create name=migration_name` - Create new migration
- `make sqlc-generate` - Generate Go code from SQL

#### Complete Setup
- `make dev-setup` - Complete development environment setup

### Docker Services

The `docker-compose.yml` provides:
- **MySQL 8.0**: Database server on port 3306
- **PhpMyAdmin**: Web interface on port 8081 (http://localhost:8081)

Default credentials:
- MySQL Root Password: `rootpassword`
- Database: `sykell_db`
- User: `sykell_user`
- Password: `sykell_password`

### Environment Variables

Copy `.env.example` to `.env` and configure:

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - MySQL connection string
- `JWT_SECRET` - Secret key for JWT tokens
- `ENVIRONMENT` - Application environment (development/production)
- `MYSQL_*` - MySQL Docker configuration

### Database Migrations

Create a new migration:
```bash
make migrate-create name=add_new_table
```

Run migrations:
```bash
make migrate-up
```

Rollback migrations:
```bash
make migrate-down
```

### SQLC Usage

1. Define your database schema in `sql/schema/`
2. Write SQL queries in `sql/queries/`
3. Generate Go code: `make sqlc-generate`
4. The generated code will be available in `internal/db/`

## Technologies Used

- [Echo](https://echo.labstack.com/) - Web framework
- [SQLC](https://sqlc.dev/) - Generate Go code from SQL
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [MySQL](https://www.mysql.com/) - Database
- [Docker](https://www.docker.com/) - Containerization
- [godotenv](https://github.com/joho/godotenv) - Environment variables
- [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - MySQL driver