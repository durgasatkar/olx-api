# --- Stage 1: Build the Go binary ---
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Cache dependencies first (improves build speeds)
COPY go.mod go.sum ./
RUN go mod download

# Copy your source code and compile the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/api ./cmd/api

# --- Stage 2: Final lightweight runtime image ---
FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy ONLY the compiled binary from the builder stage
COPY --from=builder /app/bin/api .

# Expose your dynamic port configuration
EXPOSE 8089

# Run the application using the JSON array syntax
CMD ["./api"]


