# docker/Dockerfile
FROM golang:1.24-alpine AS builder

# Install git (required for go modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy service module files first for dependency caching
COPY services/email/go.mod services/email/go.sum ./
RUN go mod download

# Copy the service source code
COPY services/email/ ./

# Build the application binary from the service module
RUN CGO_ENABLED=0 GOOS=linux go build -o email-service ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/email-service .

# Copy .env file
COPY services/email/.env .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./email-service"]
