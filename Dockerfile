# syntax=docker/dockerfile:1.4

FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install Git and libpq (for pgx etc.)
RUN apk add --no-cache git

# Cache go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build the app
WORKDIR /app/cmd/api
RUN go build -o /app/bin/app

# Final image
FROM alpine:latest

# Add certs for HTTPS and dependencies for Kafka, Postgres clients
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/app .

# Copy .env if needed at runtime (optional)
COPY .env .env

# Set environment variables
ENV GIN_MODE=release

CMD ["./app"]