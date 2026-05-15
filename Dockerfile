# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /utic-dev-server ./main.go

# Final stage
FROM alpine:3.20

WORKDIR /root/

# Copy the binary
COPY --from=builder /utic-dev-server .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy db directory for sqlc schema/queries (optional at runtime)
COPY --from=builder /app/db ./db

EXPOSE 3000

CMD ["./utic-dev-server"]
