# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Install golang-migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate sqlc code
RUN sqlc generate

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init -g cmd/server/main.go

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client stockfish

WORKDIR /root/

# Copy binary
COPY --from=builder /app/main .

# Copy migrations
COPY --from=builder /app/db/migrations ./db/migrations

# Copy Swagger docs
COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./main"]