# Stage 1: Build the Go application
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Stage 2: Final minimal image
FROM alpine:latest

# Add certificates for TLS compatibility
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the compiled binary from Stage 1
COPY --from=builder /app/main .
# Copy default config file
COPY --from=builder /app/.env.example .env

# Expose port
EXPOSE 8080

# Execute the application
CMD ["./main"]
