# Build stage
FROM golang:alpine AS builder

ENV GOTOOLCHAIN=auto

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o api-monitor ./cmd/api

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS and sqlite
RUN apk --no-cache add ca-certificates sqlite-libs

# Copy binary from builder
COPY --from=builder /app/api-monitor .
COPY --from=builder /app/internal/database/schema.sql ./internal/database/

# Expose port
EXPOSE 8080

# Run the application
CMD ["./api-monitor"]