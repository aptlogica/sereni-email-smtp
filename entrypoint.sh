#!/bin/sh

# Generate Swagger docs with dynamic host and port
echo "Generating Swagger docs with host: $SWAGGER_HOST and port: $SWAGGER_PORT"
./swag init --generalInfo cmd/server/main.go --output docs --parseDependency --parseInternal --propertyStrategy snakecase --host "$SWAGGER_HOST:$SWAGGER_PORT"

# Start the application
./email-service
