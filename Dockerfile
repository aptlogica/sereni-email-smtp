# docker/Dockerfile
FROM golang:1.24-alpine AS builder

# Install git (required for go modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy && go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o email-service ./cmd/server

# Final stage
FROM alpine:3.20

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/email-service .

# # Expose port
# EXPOSE 8080

# Run the binary
CMD ["./email-service"]
