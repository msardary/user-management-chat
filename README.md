# User Management Chat

This project is a **user management system with real-time chat functionality**, developed in **Go** using a **modular microservices architecture**. It includes **JWT-based authentication**, **role management**, **real-time messaging**, and **monitoring/logging with Prometheus, Loki, and Grafana**.

## Features

- User registration and login with JWT
- Role-based access control
- Real-time chat between users
- Dockerized deployment with Docker Compose
- Monitoring and logging with Prometheus, Loki, and Grafana
- Clean and modular project structure

## Project Structure

```
.
├── cmd/                # Application entrypoints
├── internal/           # Business logic and services
├── pkg/                # Shared utility packages
├── prometheus/         # Monitoring configurations
├── Dockerfile          # Docker image definition
├── docker-compose.yml  # Multi-container orchestration
├── env.example         # Sample environment variables
├── loki-config.yaml    # Loki logging configuration
├── promtail-config.yaml# Promtail log collector config
├── go.mod              # Go module file
└── go.sum              # Dependency checksums
```

## Requirements

- [Go 1.20+](https://golang.org/dl/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Getting Started

1. **Set up environment variables:**

   ```bash
   cp env.example .env
   # Edit the .env file with your configuration
   ```

2. **Build and run the application:**

   ```bash
   sqlc generate --file internal/db/sqlc.yaml
   docker-compose up --build
   ```

3. **Available Services:**

   - Main API: `http://localhost:8080`
   - Grafana: `http://localhost:3000` (default user/pass: `admin`)

## API Endpoints

- `POST /api/register` – Register a new user
- `POST /api/login` – Login and receive a JWT token
- `GET /api/users` – Get a list of users (requires authentication)
- `POST /api/chat/send` – Send a chat message
- `GET /api/chat/history` – Get chat message history

## Monitoring & Logging

- **Prometheus** for metric collection
- **Loki + Promtail** for log aggregation
- **Grafana** for dashboards and log visualization

## Contribution

Contributions are welcome! Please open an issue to discuss your proposal before submitting a pull request.

## License

This project is licensed under the [MIT License](LICENSE).
