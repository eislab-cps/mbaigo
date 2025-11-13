# Dockerfile for three-sensors example application

FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /build

# Copy go mod files from root
COPY go.mod go.sum ./
RUN go mod download

# Copy necessary packages
COPY pkg/ ./pkg/
COPY examples/three-sensors/ ./examples/three-sensors/

WORKDIR /build/examples/three-sensors

# Build the application
RUN go build -o app .

# Final stage - minimal runtime image
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the built application
COPY --from=builder /build/examples/three-sensors/app .

# Expose the HTTP port (from systemconfig.json)
EXPOSE 8080

# Run the application
CMD ["./app"]
