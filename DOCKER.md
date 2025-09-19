# Docker Setup for Sykell

This directory contains Docker configurations for running the complete Sykell application stack.

## Architecture

The application consists of the following services:

- **Frontend**: React application served by Nginx
- **Backend**: Go API server
- **Temporal Worker**: Go worker for background tasks
- **MySQL**: Primary database
- **Temporal**: Workflow orchestration system
- **Temporal PostgreSQL**: Database for Temporal (separate from main app DB)

## Quick Start

### Production Deployment

1. **Copy environment file:**
   ```bash
   cp .env.docker .env
   ```

2. **Edit environment variables:**
   ```bash
   nano .env
   ```
   Update `JWT_SECRET` and database passwords with secure values.

3. **Start all services:**
   ```bash
   docker-compose up -d
   ```

4. **Check service status:**
   ```bash
   docker-compose ps
   ```

5. **View logs:**
   ```bash
   docker-compose logs -f
   ```

### Development Setup

For development with hot reload:

```bash
# Start with development overrides
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

# Frontend will be available at http://localhost:5173
# Backend API at http://localhost:7070
# PHPMyAdmin at http://localhost:8080
# Temporal Web UI at http://localhost:8233
```

### Production with Optimizations

For production deployment with resource limits and security optimizations:

```bash
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
```

## Service URLs

| Service | URL | Description |
|---------|-----|-------------|
| Frontend | http://localhost | Main application |
| Backend API | http://localhost:7070 | REST API (in dev mode) |
| Temporal Web UI | http://localhost:8233 | Workflow monitoring |
| PHPMyAdmin | http://localhost:8080 | Database admin (dev only) |

## Environment Variables

Key environment variables in `.env`:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 7070 | Backend server port |
| `ENVIRONMENT` | production | Environment mode |
| `MYSQL_ROOT_PASSWORD` | rootpassword | MySQL root password |
| `MYSQL_DATABASE` | sykell_db | Database name |
| `MYSQL_USER` | sykell_user | Database user |
| `MYSQL_PASSWORD` | sykell_password | Database password |
| `JWT_SECRET` | (change this!) | JWT signing secret |
| `TEMPORAL_NAMESPACE` | default | Temporal namespace |

## Database Migration

The backend automatically runs migrations on startup. For manual migration:

```bash
# Access backend container
docker-compose exec backend sh

# Run migrations manually (if needed)
migrate -path migrations -database "mysql://user:pass@tcp(mysql:3306)/dbname" up
```

## Building Individual Services

### Backend

```bash
cd backend
docker build -t sykell-backend .
```

### Temporal Worker

```bash
cd backend
docker build -f Dockerfile.worker -t sykell-worker .
```

### Frontend

```bash
cd frontend
docker build -t sykell-frontend .
```

## Volumes and Data Persistence

The following volumes persist data:

- `mysql_data`: MySQL database files
- `temporal_postgresql_data`: Temporal PostgreSQL data

To backup data:

```bash
# Backup MySQL
docker-compose exec mysql mysqldump -u root -p sykell_db > backup.sql

# Backup Temporal
docker-compose exec temporal-postgresql pg_dump -U temporal temporal > temporal_backup.sql
```

## Scaling

Scale individual services:

```bash
# Scale workers
docker-compose up -d --scale temporal-worker=3

# Scale backend
docker-compose up -d --scale backend=2
```

Note: For multiple backend instances, you'll need a load balancer.

## Troubleshooting

### Common Issues

1. **Port conflicts**: Ensure ports 80, 3306, 7070, 7233, 8233 are available
2. **Database connection**: Wait for MySQL to be fully ready before backend starts
3. **Memory issues**: Temporal requires adequate memory (at least 2GB recommended)

### Useful Commands

```bash
# View service logs
docker-compose logs service-name

# Restart a service
docker-compose restart service-name

# Access service shell
docker-compose exec service-name sh

# Clean up everything
docker-compose down -v --remove-orphans
docker system prune -a
```

### Health Checks

All services include health checks. Check status:

```bash
docker-compose ps
```

Healthy services will show "healthy" status.

## Security Considerations

For production deployment:

1. **Change default passwords** in `.env`
2. **Use strong JWT secret**
3. **Configure HTTPS** with SSL certificates
4. **Limit exposed ports** (use prod compose file)
5. **Regular security updates** of base images
6. **Network isolation** between services
7. **Resource limits** to prevent resource exhaustion

## Monitoring and Logging

Logs are configured with rotation in production mode:
- Max file size: 10MB
- Max files: 3
- Format: JSON for structured logging

For advanced monitoring, consider adding:
- Prometheus + Grafana
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Application Performance Monitoring (APM)

## Development Notes

The development setup includes:
- Hot reload for Go services (using Air)
- Live reload for React frontend
- Volume mounts for source code
- PHPMyAdmin for database inspection
- Debug logging enabled