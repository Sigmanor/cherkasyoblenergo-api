# Cherkasyoblenergo API

[![Deploy](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml)
[![Tests](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Sigmanor/cherkasyoblenergo-api)](https://go.dev/)
[![License](https://img.shields.io/github/license/Think-Root/chappie_server)](LICENSE)
[![Releases](https://img.shields.io/github/release/Sigmanor/cherkasyoblenergo-api.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/releases)
[![Changelog](https://img.shields.io/badge/changelog-md-blue)](CHANGELOG.md)

Unofficial API service for retrieving power outage schedules from [cherkasyoblenergo.com](https://cherkasyoblenergo.com/). Get real-time and historical power outage information through a RESTful API interface.

## 📋 Table of Contents

- [Cherkasyoblenergo API](#cherkasyoblenergo-api)
  - [📋 Table of Contents](#-table-of-contents)
  - [✨ Key Features](#-key-features)
  - [🚀 Installation](#-installation)
    - [Prerequisites](#prerequisites)
    - [Setup](#setup)
  - [🔑 API Documentation](#-api-documentation)
    - [Base URL](#base-url)
    - [Available Endpoints](#available-endpoints)
  - [💻 Development](#-development)
    - [Requirements](#requirements)
    - [Local Development](#local-development)
  - [🤝 Contributing](#-contributing)
  - [❗ Troubleshooting](#-troubleshooting)
  - [🚦 Running Tests](#-running-tests)
  - [⚡ Free API Access](#-free-api-access)

## ✨ Key Features

- Real-time power outage schedule data
- Historical data access
- RESTful API interface
- Rate limiting support
- API key authentication
- Docker deployment support

## 🚀 Installation

### Prerequisites

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- PostgreSQL 17 (only if running without Docker)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/Sigmanor/cherkasyoblenergo-api.git
cd cherkasyoblenergo-api
```

2. Create `.env` file with required configurations:
```properties
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=your_strong_db_password
DB_NAME=myCoolDB
ADMIN_PASSWORD=your_strong_admin_password
SERVER_PORT=3000
```

3. Choose deployment method:

**Full Docker deployment (with PostgreSQL):**
```bash
# Create persistent volume for PostgreSQL
docker volume create postgres_data

# Deploy both app and database
docker compose -f docker-compose.app.yml -f docker-compose.db.yml up -d --build
```

**App-only deployment (for existing PostgreSQL):**
```bash
docker compose -f docker-compose.app.yml up -d --build
```

## 🔑 API Documentation

### Base URL
```
/cherkasyoblenergo/api
```

### Available Endpoints

- `POST /blackout-schedule` - Get power outage schedules
- `GET /generate-api-key` - Generate API key (admin only)
- `GET /update-api-key` - Manage API keys (admin only)

[Detailed API Documentation](API.md)

## 💻 Development

### Requirements

- Go 1.23 or higher
- PostgreSQL 17
- Docker and Docker Compose (for containerized deployment)

### Local Development

```bash
# Run locally
go run ./cmd/server/main.go

# Build
go build -o cherkasyoblenergo_api ./cmd/server/main.go
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## ❗ Troubleshooting

- **Database Connection Issues**: 
  - For Docker deployment: Check if postgres_data volume is created
  - Verify PostgreSQL credentials and connection settings in `.env`
  - For full Docker setup, ensure the db service is healthy
- **API Key Issues**: Ensure proper API key generation and rate limit configuration
- **Docker Issues**: 
  - Check Docker logs: `docker-compose logs`
  - Verify Docker network configuration
  - Ensure all required environment variables are set

## 🚦 Running Tests

To run the tests locally:
```bash
go test ./...
```

## ⚡ Free API Access

Limited free access (2 requests/minute) available for testing. Contact via [email](mailto:sigmanor@pm.me) for access.