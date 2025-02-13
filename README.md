# Cherkasyoblenergo API

[![Go Version](https://img.shields.io/github/go-mod/go-version/Sigmanor/cherkasyoblenergo-api)](https://go.dev/)
[![License](https://img.shields.io/github/license/Think-Root/chappie_server)](LICENSE)

Unofficial API service for retrieving power outage schedules from cherkasyoblenergo.com. Get real-time and historical power outage information through a RESTful API interface.

## ✨ Key Features

- Real-time power outage schedule data
- Historical data access
- RESTful API interface
- Rate limiting support
- API key authentication
- Docker deployment support

## 📋 Table of Contents

- [Cherkasyoblenergo API](#cherkasyoblenergo-api)
  - [✨ Key Features](#-key-features)
  - [📋 Table of Contents](#-table-of-contents)
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
  - [⚡ Free API Access](#-free-api-access)

## 🚀 Installation

### Prerequisites

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Setup

1. Clone the repository:
```bash
git clone https://github.com/Sigmanor/cherkasyoblenergo-api.git
cd cherkasyoblenergo-api
```

2. Create `.env` file:
```properties
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=your_strong_db_password
DB_NAME=myCoolDB
ADMIN_PASSWORD=your_strong_admin_password
SERVER_PORT=3000
```

3. Deploy:
```bash
docker-compose --env-file .env up -d --build
```

For existing PostgreSQL installations:
```bash
docker-compose -f docker-compose.app-only.yml --env-file .env up -d --build
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

- Go 1.x
- PostgreSQL

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

- **Database Connection Issues**: Verify PostgreSQL credentials and connection settings in `.env`
- **API Key Issues**: Ensure proper API key generation and rate limit configuration
- **Docker Issues**: Check Docker logs using `docker-compose logs`

## ⚡ Free API Access

Limited free access (2 requests/minute) available for testing. Contact via [email](sigmanor@pm.me) for access.

---

**Status**: Active maintenance