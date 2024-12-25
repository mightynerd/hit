# Start from the latest Golang base image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o hit .

# Use a minimal base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder image
COPY --from=builder /app/hit .

# Copy migrations
COPY --from=build /app/db/migrations ./db/migrations

# Expose port (if your application listens on a port)
EXPOSE 8080

# Command to run when starting the container
CMD ["./hit"]