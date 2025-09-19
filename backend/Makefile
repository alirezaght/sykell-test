.PHONY: build run dev test clean migrate-up migrate-down migrate-create sqlc-generate docker-up docker-down docker-logs

# Build the application
build:
	go build -o bin/server cmd/main.go

# Run the application
run: build
	./bin/server

# Run in development mode with hot reload (requires air)
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

docker-clean:
	docker-compose down -v --remove-orphans

# Database migrations (MySQL)
migrate-up:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Loading DATABASE_URL from .env file..."; \
		export $$(grep -v '^#' .env | xargs) && migrate -path migrations -database "mysql://$$DATABASE_URL" up; \
	else \
		migrate -path migrations -database "mysql://$(DATABASE_URL)" up; \
	fi

migrate-down:
	@if [ -z "$(DATABASE_URL)" ]; then \
		echo "Loading DATABASE_URL from .env file..."; \
		export $$(grep -v '^#' .env | xargs) && migrate -path migrations -database "mysql://$$DATABASE_URL" down; \
	else \
		migrate -path migrations -database "mysql://$(DATABASE_URL)" down; \
	fi

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

# Generate SQLC code
sqlc-generate:
	sqlc generate

# Install development dependencies
install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/air-verse/air@latest

# Setup development environment
setup: install-tools
	cp .env.example .env
	@echo "Don't forget to update the DATABASE_URL in .env file"

# Full development setup
dev-setup: docker-up setup
	@echo "Waiting for MySQL to be ready..."
	@sleep 10
	$(MAKE) migrate-up
	$(MAKE) sqlc-generate
	@echo "Development environment is ready!"

