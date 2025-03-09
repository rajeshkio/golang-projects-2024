# Use a small Go base image
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies and build the Go binary
RUN go mod tidy && go build -o pixiedocs main.go

# Use a lightweight image for the final container
FROM alpine:latest

# Install necessary dependencies (if any)
RUN apk add --no-cache ca-certificates poppler-utils

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/pixiedocs /app/pixiedocs

# Expose the port (if your app runs as a web service)
EXPOSE 8080

# Set the entrypoint to run the binary
ENTRYPOINT ["/app/pixiedocs"]
