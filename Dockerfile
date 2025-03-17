FROM golang:1.22-alpine

WORKDIR /mnt

# Install necessary tools
RUN apk update && apk add --no-cache git ca-certificates tzdata postgresql-client && \
    update-ca-certificates


# Copy the entire project into /mnt
COPY . .

# Download go modules.
RUN go mod download

# The Dockerfile doesn't build anything here.  The build happens
# when build.sh is run *inside* the container.

EXPOSE 9090