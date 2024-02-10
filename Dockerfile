# Use a minimal base image for efficiency
FROM golang:alpine AS builder

# Set the working directory within the container
WORKDIR /app

# Copy Go source code and Fiber dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy all other application files and folders
COPY . /app

# Build the Go application (replace "main" with your entrypoint file name)
RUN go build -o main

# Create a new runtime image based on the builder image
FROM alpine

# Set the working directory
WORKDIR /app

# Copy the built binary
COPY --from=builder /app/main .
COPY .env /app

# Expose the port your Fiber app listens on (e.g., 8080)
EXPOSE 3333

# Set the command to run the application when the container starts
CMD ["/app/main"]