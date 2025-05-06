#!/bin/bash

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Swag not found, installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Generate Swagger documentation
echo "Generating Swagger documentation..."
swag init -g ./cmd/server/main_swagger.go -o ./docs

echo "Swagger documentation generated successfully!"
echo "You can access the Swagger UI at http://localhost:8443/swagger/index.html when the server is running."
