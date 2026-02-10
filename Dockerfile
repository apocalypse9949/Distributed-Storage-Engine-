# Use the official Golang image as the base image
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to download dependencies
COPY go.mod ./
COPY go.sum ./

# Download the Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# The CGO_ENABLED=0 flag disables cgo, which makes the binary statically linked
# This is useful for creating small, self-contained Docker images
RUN CGO_ENABLED=0 go build -o /palantir cmd/palantir/main.go

# Use a minimal Alpine Linux image for the final stage
# This creates a very small final image, as it only contains the compiled binary
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /palantir .

# Expose any necessary ports (e.g., for gRPC or HTTP)
# We will add actual ports later when we implement the API
EXPOSE 50051

# Command to run the executable
CMD ["./palantir"]
