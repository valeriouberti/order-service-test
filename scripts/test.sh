#!/bin/sh
set -e

echo "Running tests..."
go test -v ./... -cover

# If you want to generate a coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo "Tests completed!"