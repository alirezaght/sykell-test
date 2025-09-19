# Backend Dockerfile
# Multi-stage build for Go application

# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests and tzdata for timezone
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/server .

# Copy migrations and other necessary files
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/sql ./sql

# Create a non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Change ownership of the app directory
RUN chown -R appuser:appgroup /root

# Switch to non-root user
USER appuser

# Expose port (should match your application port)
EXPOSE 7070

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:7070/health || exit 1

# Command to run the application
CMD ["./server"]