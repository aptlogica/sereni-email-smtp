# sereni-email-smtp Makefile

# Variables
BINARY_NAME=sereni-email-smtp
VERSION?=$(shell git describe --tags --always --dirty)
COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildTime=${BUILD_TIME}"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
COVER_DIR=coverage
COVER_PROFILE=$(COVER_DIR)/coverage.out
COVER_HTML=$(COVER_DIR)/coverage.html

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

.PHONY: help build test test-coverage coverage coverage-func clean run deps lint format security docker-build docker-run

# Default target
all: deps format lint test build

# Help target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@echo '  build            - Build the application'
	@echo '  run              - Build and run the application'
	@echo '  test             - Run all tests'
	@echo '  test-coverage    - Run tests and show coverage'
	@echo '  coverage         - Alias for test-coverage'
	@echo '  coverage-func    - Show coverage by function'
	@echo '  clean            - Clean build artifacts'
	@echo '  deps             - Download and install dependencies'
	@echo '  lint             - Run golangci-lint'
	@echo '  lint-fix         - Run golangci-lint with auto-fix'
	@echo '  format           - Format Go code'
	@echo '  security         - Run security scan'

# Build the application
build: ## Build the application
	@echo "${GREEN}Building ${BINARY_NAME}...${NC}"
	$(GOBUILD) ${LDFLAGS} -o bin/${BINARY_NAME} ./cmd/server
	@echo "${GREEN}Build completed successfully!${NC}"

# Run tests
test: ## Run all tests
	@echo "${GREEN}Running tests...${NC}"
	@mkdir -p $(COVER_DIR)
	$(GOTEST) -v -race -coverprofile=$(COVER_PROFILE) -covermode=atomic ./...
	@echo "${GREEN}Tests completed!${NC}"

# Run tests with coverage
test-coverage: test ## Run tests and show coverage
	@echo "${GREEN}Generating coverage report...${NC}"
	$(GOCMD) tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)
	@echo "${BLUE}Coverage report generated: $(COVER_HTML)${NC}"

coverage: test-coverage ## Alias for test-coverage

coverage-func: ## Show coverage by function
	$(GOCMD) tool cover -func=$(COVER_PROFILE)

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "${YELLOW}Cleaning...${NC}"
	$(GOCLEAN)
	rm -rf bin/
	rm -rf $(COVER_DIR)
	@echo "${GREEN}Clean completed!${NC}"

# Run the application
run: build ## Build and run the application
	@echo "${GREEN}Running ${BINARY_NAME}...${NC}"
	./bin/${BINARY_NAME}

# Install dependencies
deps: ## Download and install dependencies
	@echo "${GREEN}Installing dependencies...${NC}"
	$(GOMOD) tidy
	$(GOMOD) download
	@echo "${GREEN}Dependencies installed!${NC}"

# Run linting
lint: ## Run golangci-lint
	@echo "${GREEN}Running linter...${NC}"
	@which golangci-lint > /dev/null || (echo "${RED}golangci-lint not found. Install with: 'go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest'${NC}" && exit 1)
	golangci-lint run
	@echo "${GREEN}Linting completed!${NC}"

# Format code
format: ## Format Go code
	@echo "${GREEN}Formatting code...${NC}"
	$(GOFMT) -s -w .
	@echo "${GREEN}Code formatted!${NC}"

# Security scan
security: ## Run security scan with gosec
	@echo "${GREEN}Running security scan...${NC}"
	@which gosec > /dev/null || (echo "${RED}gosec not found. Install with: 'go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest'${NC}" && exit 1)
	gosec ./...
	@echo "${GREEN}Security scan completed!${NC}"

# Docker build
docker-build: ## Build Docker image
	@echo "${GREEN}Building Docker image...${NC}"
	docker build -t ${BINARY_NAME}:${VERSION} .
	docker tag ${BINARY_NAME}:${VERSION} ${BINARY_NAME}:latest
	@echo "${GREEN}Docker image built: ${BINARY_NAME}:${VERSION}${NC}"

# Docker run
docker-run: ## Run Docker container with SMTP testing
	@echo "${GREEN}Running Docker container...${NC}"
	docker run --rm -p 8080:8080 \
		-e SMTP_HOST=${SMTP_HOST} \
		-e SMTP_PORT=${SMTP_PORT} \
		-e SMTP_USERNAME=${SMTP_USERNAME} \
		-e SMTP_PASSWORD=${SMTP_PASSWORD} \
		-e FROM_EMAIL=${FROM_EMAIL} \
		${BINARY_NAME}:latest

# Start development environment
dev-env: ## Start development environment with Docker Compose
	@echo "${GREEN}Starting development environment...${NC}"
	docker-compose up -d redis mailhog
	@echo "${BLUE}Services starting. MailHog UI available at http://localhost:8025${NC}"

# Stop development environment
dev-env-stop: ## Stop development environment
	@echo "${YELLOW}Stopping development environment...${NC}"
	docker-compose down

# Send test email (requires running server)
test-email: ## Send a test email
	@echo "${GREEN}Sending test email...${NC}"
	curl -X POST http://localhost:8080/send \
		-H "Content-Type: application/json" \
		-d '{"to":["test@example.com"],"subject":"Test Email","body":"This is a test email from sereni-email-smtp"}'

# Generate Swagger docs
swagger: ## Generate swagger documentation
	@echo "${GREEN}Generating Swagger documentation...${NC}"
	@which swag > /dev/null || (echo "${RED}swag not found. Install with: 'go install github.com/swaggo/swag/cmd/swag@latest'${NC}" && exit 1)
	swag init -g cmd/server/main.go
	@echo "${GREEN}Swagger docs generated!${NC}"

# Pre-commit checks
pre-commit: format lint test security ## Run all pre-commit checks
	@echo "${GREEN}All pre-commit checks passed!${NC}"
