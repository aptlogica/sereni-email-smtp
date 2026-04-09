
## docker/Dockerfile
FROM golang:1.26.2-alpine@sha256:c2a1f7b2095d046ae14b286b18413a05bb82c9bca9b25fe7ff5efef0f0826166 AS builder
RUN go version

# Install git (required for go modules)
RUN apk add --no-cache git

# Install swag CLI (pinned to v1.16.4 - commit 0b9e347c196710ea155a147782bf51707a600c2c)
RUN git clone https://github.com/swaggo/swag.git /tmp/swag && \
    cd /tmp/swag && \
    git checkout 0b9e347c196710ea155a147782bf51707a600c2c && \
    go install ./cmd/swag && \
    rm -rf /tmp/swag

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy && go mod download

# Copy source code
COPY . .


# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sereni-email-smtp ./cmd/server


# Copy swag binary for later use
RUN cp /go/bin/swag /app/swag


# Final stage
FROM alpine:3.20@sha256:a4f4213abb84c497377b8544c81b3564f313746700372ec4fe84653e4fb03805

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/sereni-email-smtp .
COPY --from=builder /app/swag .



# # Expose port
# EXPOSE 8080




ENTRYPOINT ["./sereni-email-smtp"]
