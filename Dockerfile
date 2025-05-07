# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# RUN go install github.com/kyleconroy/sqlc/cmd/sqlc@v1.16.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# COPY internal/db/sqlc.yaml /app/sqlc.yaml

# RUN sqlc version

# RUN sqlc generate

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/app ./cmd/app

# Run stage
FROM debian:bookworm-slim

WORKDIR /app

# Install minimal dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app /app/app
COPY .env .env

# Create logs directory
RUN mkdir -p /logs && chmod 777 /logs

EXPOSE 8080

CMD ["/app/app"]