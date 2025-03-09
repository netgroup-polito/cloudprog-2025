# Use the official Golang image to create a build artifact.
FROM golang:1.24 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) go build -ldflags="-s -w" -o main ./cmd/main.go

# Start a new stage from scratch
FROM alpine:3.21.2  

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main /app/main

# Copy the openapi.yml file from the previous stage
COPY ./openapi.yml /app/openapi.yml

# Command to run the executable
ENTRYPOINT ["/app/main"]
