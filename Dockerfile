# Use Go 1.22 as base image
FROM golang:1.22-alpine AS builder

# Install dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o realentity-node cmd/main.go

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/realentity-node .

# Copy default config
COPY --from=builder /app/config.json .

# Expose the default port
EXPOSE 4001

# Command to run
CMD ["./realentity-node"]
