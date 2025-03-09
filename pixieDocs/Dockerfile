# Build stage
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git build-base

# Set working directory
WORKDIR /build

# Copy go.mod and go.sum files first for better caching
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations for container environment
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -X main.version=$(git describe --tags --always) -X main.buildTime=$(date +%Y-%m-%dT%H:%M:%S)" \
    -o pixiedocs .

# Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache qpdf ca-certificates tzdata && \
    update-ca-certificates

# Create a non-root user to run the application
RUN adduser -D -h /app appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /build/pixiedocs /app/
COPY --from=builder /build/static /app/static

# Set ownership for security
RUN chown -R appuser:appuser /app

# Use non-root user for security
USER appuser

# Expose application port
EXPOSE 8080

# Configure health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:8080/health || exit 1

# Set environment variables
ENV GIN_MODE=release
ENV GOTENBERG_URL=https://pixiedocs.rajesh-kumar.in/

# Run the application
ENTRYPOINT ["/app/pixiedocs"]