#!/bin/bash

echo "🚀 Starting Sykell Development Environment"

# Start MySQL first
echo "📊 Starting MySQL..."
docker-compose up -d mysql

# Wait for MySQL to be ready
echo "⏳ Waiting for MySQL to be ready..."
docker-compose exec mysql mysqladmin ping -h"localhost" --silent

# Run migrations
echo "🔄 Running database migrations..."
docker-compose --profile migration up migrate

# Start all other services
echo "🌟 Starting all services..."
docker-compose up -d

echo "✅ All services started!"
echo "Frontend: http://localhost"
echo "Backend API: http://localhost:7070"
echo "Temporal UI: http://localhost:8233"

# Show service status
docker-compose ps