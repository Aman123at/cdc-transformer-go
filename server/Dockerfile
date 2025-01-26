# Stage 1: Build the Go application
FROM golang:1.23.4-alpine AS builder

# Install system dependencies (if needed)
# RUN apk add --no-cache git gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cdc-go ./main.go

# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install dependencies (if needed, e.g., for SSL certificates)
# RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/cdc-go /app/cdc-go

# Copy the .env file
COPY .env .env

COPY server.crt server.crt
COPY server.key server.key

# Expose the port your app runs on
EXPOSE 8000

# Run the application
CMD ["/app/cdc-go"]