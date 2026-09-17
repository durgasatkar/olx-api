# --- Stage 1: Build the Go binaries ---
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# 1. Build the main API binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api ./cmd/api

# 2. Build the migration binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/migrate ./cmd/migrate

# --- Stage 2: Final lightweight runtime image ---
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy both compiled binaries from the builder stage
COPY --from=builder /app/bin/api .
COPY --from=builder /app/bin/migrate .

# Copy your migrations folder containing your SQL files
COPY --from=builder /app/migrations ./migrations

EXPOSE 8089
CMD ["./api"]


