# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.io,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

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