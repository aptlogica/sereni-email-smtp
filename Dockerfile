# Copy entrypoint script
COPY entrypoint.sh .

# Make entrypoint executable
RUN chmod +x entrypoint.sh
# docker/Dockerfile
FROM golang:1.24.4-alpine AS builder

# Install git (required for go modules)
RUN apk add --no-cache git

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o email-service cmd/server/main.go

# Copy swag binary for later use
RUN cp /go/bin/swag /app/swag


# Final stage
FROM alpine:3.20

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/email-service .
COPY --from=builder /app/swag .


# Expose port
EXPOSE 8082


# Set default env vars (can be overridden)
ENV SWAGGER_HOST=localhost
ENV SWAGGER_PORT=8082

ENTRYPOINT ["./entrypoint.sh"]
