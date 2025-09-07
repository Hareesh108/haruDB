# Use latest Go (1.24) for build
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

# Build harudb binary
RUN go build -o harudb ./cmd/server

# Run in a lightweight container
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/harudb /usr/local/bin/harudb

EXPOSE 54321
CMD ["harudb"]
