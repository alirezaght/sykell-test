# Migration Dockerfile
FROM migrate/migrate:latest

# Copy migration files
COPY backend/migrations /migrations

# Set the working directory
WORKDIR /migrations

# Default command to run migrations up
# The DATABASE_URL will be passed as environment variable
ENTRYPOINT ["sh", "-c", "migrate -path /migrations -database \"$DATABASE_URL\" up"]