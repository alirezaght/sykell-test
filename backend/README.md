# Sykell Backend

Go backend server for the Sykell web crawler application, featuring a clean architecture with SOLID principles, Temporal for reliable task processing, and MySQL for data persistence.

## Technology Stack

- **Go 1.21+** - Programming language
- **Echo** - Web framework for HTTP server
- **MySQL 8.0** - Database
- **SQLC** - Type-safe SQL code generation
- **Temporal** - Durable task execution
- **Zap** - Structured logging
- **JWT** - Authentication
- **golang-migrate** - Database migrations
- **Air** - Hot reload for development

## Prerequisites

- **Go 1.21+**
- **Docker & Docker Compose** (for local development)
- **Make** (for build automation)

## Quick Start

### 1. Clone and Setup

```bash
cd backend/

# Install development tools and setup environment
make setup

# This will:
# - Install sqlc, migrate, air tools
# - Copy .env.example to .env
```

### 2. Environment Configuration

Edit the `.env` file with your configuration:

```bash
# Server Configuration
PORT=7070
ENVIRONMENT=development

# Database Configuration (MySQL)
DATABASE_URL=sykell_user:sykell_password@tcp(localhost:3306)/sykell_db?charset=utf8mb4&parseTime=True&loc=Local

# MySQL Docker Configuration
MYSQL_ROOT_PASSWORD=rootpassword
MYSQL_DATABASE=sykell_db
MYSQL_USER=sykell_user
MYSQL_PASSWORD=sykell_password

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Temporal Configuration
TEMPORAL_HOST_PORT=localhost:7233
TEMPORAL_NAMESPACE=default
BACKEND_URL=http://localhost:7070
```

### 3. Start Development Environment

```bash
# Start all services (MySQL, Temporal, Workers)
make dev-setup

# This will:
# - Start Docker services
# - Wait for MySQL to be ready
# - Run database migrations
# - Generate SQLC code
```

### 4. Start Development Server

```bash
# Start with hot reload
make dev

# Or build and run manually
make run
```

The API will be available at **http://localhost:7070**

### Troubleshooting
- Dev (Air) error: `no Go files in .../backend`
    - Ensure you run commands from the `backend/` directory.
    - We force Air to use `.air.toml` (builds `./cmd/main.go`). If you still see this, reinstall Air and try again:
        ```bash
        make install-tools
        make dev
        ```
    - If a stale tmp binary path is referenced: remove `backend/tmp/` and retry.
        ```bash
        rm -rf tmp
        make dev
        ```
- Container name conflict (`/sykell_mysql` already in use):
   ```
   Error response from daemon: Conflict. The container name "/sykell_mysql" is already in use
   ```
   The root stack and `backend/docker-compose.yml` both define a MySQL service named `sykell_mysql`. Don’t run both stacks at the same time. Fix it with one of these 

   ```bash
   # Stop the other stack, then start the one you need
   # From repo root (if root stack is running)
   docker compose down -v
   # From backend/ (if backend stack is running)
   docker compose down -v
   
   ```
   

## Make Commands

| Command | Description |
|---------|-------------|
| `make setup` | Install dev tools and create .env file |
| `make dev-setup` | Full development environment setup |
| `make dev` | Start development server with hot reload |
| `make build` | Build the application binary |
| `make run` | Build and run the application |
| `make test` | Run all tests |
| `make clean` | Clean build artifacts |


### Docker Commands

| Command | Description |
|---------|-------------|
| `make docker-up` | Start Docker services |
| `make docker-down` | Stop Docker services |
| `make docker-logs` | View Docker logs |
| `make docker-clean` | Stop and remove all containers/volumes |

### Database Commands

| Command | Description |
|---------|-------------|
| `make migrate-up` | Run database migrations |
| `make migrate-down` | Rollback database migrations |
| `make migrate-create name=migration_name` | Create new migration |
| `make sqlc-generate` | Generate Go code from SQL queries |

### Tool Installation

| Command | Description |
|---------|-------------|
| `make install-tools` | Install sqlc, migrate, air tools |

## Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | HTTP server port | `7070` | Yes |
| `ENVIRONMENT` | Environment mode | `development` | Yes |
| `DATABASE_URL` | MySQL connection string | - | Yes |
| `JWT_SECRET` | JWT signing secret | - | Yes |
| `TEMPORAL_HOST_PORT` | Temporal server address | `localhost:7233` | Yes |
| `TEMPORAL_NAMESPACE` | Temporal namespace | `default` | Yes |
| `BACKEND_URL` | Backend base URL for workers | `http://localhost:7070` | Yes |

## Architecture & SOLID Principles

### Feature-First Structure

Each feature is organized in its own package following SOLID principles:

```
internal/
├── config/              # Configuration management
├── middleware/          # HTTP middleware (JWT, logging)
├── db/                  # Generated SQLC code
├── utils/               # Shared utilities
├── temporal/            # Temporal service integration
└── [feature]/           # Feature packages
    ├── handler.go       # HTTP handlers (network layer)
    ├── service.go       # Business logic
    ├── repo.go         # Data access layer (interface)
    ├── dto.go          # Data Transfer Objects
    ├── [feature]_test.go # Unit tests
    └── [extensions].go  # Feature extensions
```

### Layer Responsibilities

#### 1. Handler Layer (`*_handler.go`)
- **Responsibility**: HTTP request/response handling
- **Concerns**: Request validation, response formatting, authentication
- **Dependencies**: Service layer only

```go
// Example: handler only handles HTTP concerns
func (h *Handler) CreateUser(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }
    
    // Delegate to service layer
    user, err := h.service.CreateUser(c.Request().Context(), req)
    if err != nil {
        return handleServiceError(err)
    }
    
    return c.JSON(http.StatusCreated, user)
}
```

#### 2. Service Layer (`*_service.go`)
- **Responsibility**: Business logic and orchestration
- **Concerns**: Business rules, data transformation, external service coordination
- **Dependencies**: Repository interfaces, external services

```go
// Example: service handles business logic
func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
    // Business validation
    if err := s.validateUserData(req); err != nil {
        return nil, err
    }
    
    // Hash password (business logic)
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }
    
    // Delegate to repository
    return s.repo.Create(ctx, req.Email, hashedPassword)
}
```

#### 3. Repository Layer (`*_repo.go`)
- **Responsibility**: Data access abstraction
- **Concerns**: Database operations, query composition
- **Dependencies**: Database connection only

```go
// Repository is always an interface
type Repo interface {
    Create(ctx context.Context, email string, passwordHash string) (*UserResponse, error)
    GetByEmail(ctx context.Context, email string) (*UserResponse, error)
    GetByID(ctx context.Context, id string) (*UserResponse, error)
}

// Implementation uses SQLC generated code
type userRepo struct {
    queries *db.Queries
}

func (r *userRepo) Create(ctx context.Context, email string, passwordHash string) (*UserResponse, error) {
    // Use SQLC generated methods
    user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
        Email:        email,
        PasswordHash: passwordHash,
    })
    // ... error handling and response mapping
}
```

#### 4. DTO Layer (`*_dto.go`)
- **Responsibility**: Data contracts and validation
- **Concerns**: Request/response structures, validation rules

```go
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
}
```

## Working with SQLC

### Adding New Queries

1. **Write SQL query** in `sql/queries/[feature].sql`:

```sql
-- name: GetUserByEmail :one
SELECT id, email, password_hash, created_at, updated_at
FROM users
WHERE email = ? LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (id, email, password_hash)
VALUES (?, ?, ?)
RETURNING id, email, created_at, updated_at;
```

2. **Generate Go code**:

```bash
make sqlc-generate
```

3. **Use in repository**:

```go
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*UserResponse, error) {
    user, err := r.queries.GetUserByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    
    return &UserResponse{
        ID:    user.ID,
        Email: user.Email,
    }, nil
}
```

### SQLC Configuration

The `sqlc.yaml` file configures code generation:

```yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "./sql/queries"    # SQL queries location
    schema: "./migrations"      # Migration files for schema
    gen:
      go:
        package: "db"
        out: "./internal/db"    # Generated code location
        emit_interface: true    # Generate interfaces
        emit_json_tags: true    # Add JSON tags
```

## Adding a New Feature

Follow these steps to add a new feature that adheres to SOLID principles:

### 1. Create Feature Package

```bash
mkdir internal/newfeature
```

### 2. Define Data Structures (`dto.go`)

```go
package newfeature

// Request/Response structures
type CreateRequest struct {
    Name string `json:"name" validate:"required"`
}

type Response struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

### 3. Create Repository Interface (`repo.go`)

```go
package newfeature

import "context"

// Always define as interface for testability
type Repo interface {
    Create(ctx context.Context, name string) (*Response, error)
    GetByID(ctx context.Context, id string) (*Response, error)
    List(ctx context.Context) ([]Response, error)
}

// Implementation
type repo struct {
    queries *db.Queries
}

func NewRepo(queries *db.Queries) Repo {
    return &repo{queries: queries}
}
```

### 4. Implement Service (`service.go`)

```go
package newfeature

import "context"

type Service interface {
    Create(ctx context.Context, req CreateRequest) (*Response, error)
    Get(ctx context.Context, id string) (*Response, error)
}

type service struct {
    repo Repo
    // other dependencies (logger, external services)
}

func NewService(repo Repo) Service {
    return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateRequest) (*Response, error) {
    // Business validation
    if err := s.validateRequest(req); err != nil {
        return nil, err
    }
    
    // Business logic
    // ...
    
    // Delegate to repository
    return s.repo.Create(ctx, req.Name)
}
```

### 5. Create HTTP Handlers (`handler.go`)

```go
package newfeature

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

type Handler struct {
    service Service
}

func NewHandler(service Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Create(c echo.Context) error {
    var req CreateRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }
    
    result, err := h.service.Create(c.Request().Context(), req)
    if err != nil {
        return handleError(err)
    }
    
    return c.JSON(http.StatusCreated, result)
}

// Register routes
func (h *Handler) RegisterRoutes(g *echo.Group) {
    g.POST("/newfeature", h.Create)
    g.GET("/newfeature/:id", h.Get)
}
```

### 6. Add Database Migration

```bash
make migrate-create name=create_newfeature_table
```

Edit the generated migration files:

```sql
-- migrations/000X_create_newfeature_table.up.sql
CREATE TABLE newfeature (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 7. Add SQL Queries

Create `sql/queries/newfeature.sql`:

```sql
-- name: CreateNewFeature :one
INSERT INTO newfeature (id, name)
VALUES (?, ?)
RETURNING id, name, created_at, updated_at;

-- name: GetNewFeatureByID :one
SELECT id, name, created_at, updated_at
FROM newfeature
WHERE id = ? LIMIT 1;
```

### 8. Generate and Wire Up

```bash
# Generate SQLC code
make sqlc-generate

# Run migration
make migrate-up
```

### 9. Register in Main Application

In `cmd/main.go`:

```go
// Initialize dependencies
newFeatureRepo := newfeature.NewRepo(queries)
newFeatureService := newfeature.NewService(newFeatureRepo)
newFeatureHandler := newfeature.NewHandler(newFeatureService)

// Register routes
apiV1 := e.Group("/api/v1")
newFeatureHandler.RegisterRoutes(apiV1)
```

### 10. Add Tests

Create `newfeature_test.go`:

```go
package newfeature

import (
    "testing"
    "context"
)

func TestService_Create(t *testing.T) {
    // Mock repository
    mockRepo := &mockRepo{}
    service := NewService(mockRepo)
    
    // Test business logic
    req := CreateRequest{Name: "test"}
    result, err := service.Create(context.Background(), req)
    
    // Assertions
    assert.NoError(t, err)
    assert.Equal(t, "test", result.Name)
}
```

## Testing Strategy

### Unit Tests
- **Location**: Same package as source code (`*_test.go`)
- **Focus**: Business logic in service layer
- **Mocking**: Repository interfaces

### Integration Tests
- **Location**: `tests/integration/`
- **Focus**: Database interactions, API endpoints
- **Setup**: Test database, real dependencies

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/user -v

# With coverage
go test -cover ./...
```

## Logging

The application uses Zap for structured logging:

```go
import "go.uber.org/zap"

// In service layer
func (s *service) CreateUser(ctx context.Context, req CreateUserRequest) error {
    s.logger.Info("Creating user",
        zap.String("email", req.Email),
        zap.String("request_id", getRequestID(ctx)),
    )
    
    // ... business logic
    
    s.logger.Error("Failed to create user",
        zap.Error(err),
        zap.String("email", req.Email),
    )
}
```

## Authentication & Authorization

### JWT Authentication
- **Bearer tokens** for API endpoints
- **Cookie-based auth** for SSE connections (header limitations)

### Middleware Usage

```go
// Protected routes
protected := apiV1.Group("")
protected.Use(middleware.JWTAuth(jwtSecret))
protected.GET("/profile", userHandler.GetProfile)

// Public routes
apiV1.POST("/auth/login", authHandler.Login)
apiV1.POST("/auth/register", authHandler.Register)
```

## Temporal Integration

### Worker Implementation

Workers are implemented in the feature packages:

```go
// internal/crawl/worker.go
func (w *CrawlWorker) CrawlActivity(ctx context.Context, params CrawlParams) error {
    // Long-running crawl logic
    // Fault-tolerant and retryable
}

// Register in cmd/worker/main.go
worker.RegisterActivity(crawlWorker.CrawlActivity)
```

### Task Scheduling

```go
// In service layer
func (s *service) StartCrawl(ctx context.Context, urlID string) error {
    // Queue task to Temporal
    workflowOptions := client.StartWorkflowOptions{
        ID:        fmt.Sprintf("crawl-%s", urlID),
        TaskQueue: "crawl-queue",
    }
    
    _, err := s.temporalClient.ExecuteWorkflow(ctx, workflowOptions, 
        "CrawlWorkflow", CrawlParams{URLID: urlID})
    return err
}
```