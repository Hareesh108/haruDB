# =========================
# 1️⃣ Build stage
# =========================
FROM golang:1.24-alpine AS builder

# Install git for go modules
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod and go.sum for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build server and CLI binaries
RUN go build -o harudb ./cmd/server
RUN go build -o haru-cli ./cmd/cli

# =========================
# 2️⃣ Runtime stage
# =========================
FROM alpine:latest

# Add non-root user
RUN addgroup -S harudb && adduser -S -G harudb harudb

WORKDIR /app

# Copy binaries
COPY --from=builder /app/harudb /usr/local/bin/harudb
COPY --from=builder /app/haru-cli /usr/local/bin/haru-cli

# Create data directory
RUN mkdir -p /app/data && chown -R harudb:harudb /app/data

USER harudb

# Expose HaruDB server port
EXPOSE 54321

# Set default command to server
CMD ["harudb", "--data-dir", "/app/data"]
