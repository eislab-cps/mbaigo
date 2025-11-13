# Multi-stage Dockerfile for mbaigo
# Used for both core systems (ESR, Orchestrator) and application systems

FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the mbaigo CLI
RUN make build

# Final stage - minimal runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the built binary
COPY --from=builder /build/bin/mbaigo /usr/local/bin/mbaigo

# Create directories for core systems configs
RUN mkdir -p /app/systems/esr /app/systems/orchestrator

# Set working directory
WORKDIR /app

# Default command (can be overridden in docker-compose)
CMD ["mbaigo", "--help"]
