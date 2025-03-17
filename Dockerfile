FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install required tools
RUN apk update && \
    apk add --no-cache git gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/order-service-test ./cmd/api

# Create final lightweight image
FROM alpine:3.18

WORKDIR /app

# Install necessary packages
RUN apk update && \
    apk add --no-cache ca-certificates tzdata postgresql-client && \
    update-ca-certificates

# Copy the binary from builder
COPY --from=builder /app/order-service-test /app/order-service-test

# Copy scripts and migrations
COPY ./scripts/ /app/scripts/
COPY ./migrations/ /app/migrations/

# Make scripts executable
RUN chmod +x /app/scripts/*.sh

# Default port
EXPOSE 9090

# Run the application
CMD ["/app/order-service-test"]