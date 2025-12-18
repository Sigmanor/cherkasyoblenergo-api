<div align="center">

# Cherkasyoblenergo API

[![EN](https://img.shields.io/badge/English-0e7837.svg)](README.md) [![UA](https://img.shields.io/badge/Ukrainian-c9c9c9.svg)](README_UA.md)

[![Deploy](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml)
[![Tests](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Sigmanor/cherkasyoblenergo-api)](https://go.dev/)
[![License](https://img.shields.io/github/license/Think-Root/chappie_server)](LICENSE)
[![Releases](https://img.shields.io/github/release/Sigmanor/cherkasyoblenergo-api.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/releases)
[![Changelog](https://img.shields.io/badge/changelog-md-blue)](CHANGELOG.md)

Unofficial API service for retrieving power outage schedules from [cherkasyoblenergo.com](https://cherkasyoblenergo.com/). Get real-time and historical power outage information through a RESTful API interface.

</div>

## ✨ Key Features

- Real-time power outage schedule data
- Historical data access
- RESTful API interface
- Rate limiting support
- API key authentication
- Webhook notifications

## 🚀 Installation

### Prerequisites

- Go 1.24 or higher
- PostgreSQL 17

### Setup

1. **Install PostgreSQL 17**

   Follow the [official PostgreSQL installation guide](https://www.postgresql.org/download/) to install PostgreSQL on your system.

2. **Clone the repository**

   ```bash
   git clone https://github.com/Sigmanor/cherkasyoblenergo-api.git
   cd cherkasyoblenergo-api
   ```

3. **Configure environment variables**

   Create a `.env` file in the root directory:

   ```properties
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=root
   DB_PASSWORD=your_strong_db_password
   DB_NAME=myCoolDB
   ADMIN_PASSWORD=your_strong_admin_password
   SERVER_PORT=3000
   ```

4. **Run the application**

   ```bash
   go run ./cmd/server/main.go
   ```

   The application will automatically create the required database on first run.

### Building

To build the application for production:

```bash
go build -o cherkasyoblenergo_api ./cmd/server/main.go
```

## 🔑 API Documentation

### Base URL

```
/cherkasyoblenergo/api
```

### Available Endpoints

- `GET /blackout-schedule` - Get power outage schedules
- `POST /api-keys` - Create API key (admin only)
- `PATCH /api-keys` - Rotate key or update rate limit (admin only)
- `DELETE /api-keys` - Delete API key (admin only)
- `POST /webhook` - Register webhook URL
- `DELETE /webhook` - Delete webhook
- `GET /webhook` - Get webhook status

[Detailed API Documentation](API.md)

### Webhook Notifications

The API supports webhook notifications that automatically send new power outage schedules to your endpoint whenever they become available:

- **Automatic delivery** - Get notified immediately when new schedules are parsed
- **Retry logic** - Failed deliveries are retried with exponential backoff
- **Automatic disabling** - Webhooks are disabled after 3 consecutive failures
- **Secure headers** - Includes authentication token for verification

Register your webhook URL via `POST /webhook` and receive real-time updates without polling.

## 🚦 Running Tests

To run the tests locally:

```bash
go test ./...
```

## ⚡ Free API Access

Since I host the app for my own needs, I can provide you with limited access (6 req/min) to the API and webhook notifications for free. Contact via [email](mailto:dock-brunt-rarity@duck.com) for access.
