# Build stage
FROM golang:1.24.3-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with CGO enabled
ENV CGO_ENABLED=1
RUN go build -o main ./cmd/students-api/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Create storage directory
RUN mkdir -p /app/storage && chmod 777 /app/storage

# Copy the binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/config/production.yaml /app/config/production.yaml

# Set environment variables
ENV CONFIG_PATH=/app/config/production.yaml

EXPOSE 4040

# Run the application
CMD ["./main"] 

