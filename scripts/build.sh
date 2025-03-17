#!/bin/sh
set -e

echo "Building the application..."
go build -o order-service-test ./cmd/api/main.go
echo "Build completed successfully!"