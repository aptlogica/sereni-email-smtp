
## docker/Dockerfile
FROM golang:1.27.3-alpine3.24@sha256:fa3d596898d9f126af7c01d14ee7597bab3f39cccf392ef36b69fe77cd681f18 AS builder
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
FROM alpine:3.24@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/sereni-email-smtp .
COPY --from=builder /app/swag .



# # Expose port
# EXPOSE 8080




ENTRYPOINT ["./sereni-email-smtp"]
