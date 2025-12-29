<div align="center">

# Cherkasyoblenergo API

[![EN](https://img.shields.io/badge/English-0e7837.svg)](README.md) [![UA](https://img.shields.io/badge/Ukrainian-c9c9c9.svg)](README_UA.md)

[![Deploy](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml)
[![Tests](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Sigmanor/cherkasyoblenergo-api)](https://go.dev/)
[![License](https://img.shields.io/github/license/Think-Root/chappie_server)](LICENSE)
[![Releases](https://img.shields.io/github/release/Sigmanor/cherkasyoblenergo-api.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/releases)
[![Changelog](https://img.shields.io/badge/changelog-md-blue)](CHANGELOG.md)

**Unofficial API** for retrieving power outage schedules from [cherkasyoblenergo.com](https://cherkasyoblenergo.com/)

</div>

## Public Instance

A public instance of this API is available at:

```
https://hue.pp.ua/cherkasyoblenergo/api/
```

Feel free to use it for your projects. However, please note that this instance is provided "as is" without any guarantees of availability or uptime. For production use cases or if you need reliability, consider self-hosting your own instance.

> [!NOTE]
> This public instance has a rate limit of 30 requests per minute per IP address.

## Features

- Real-time power outage schedule data
- Historical data access
- RESTful API interface
- IP-based rate limiting
- Response caching
- Optional API key authentication

## Installation

### Prerequisites

- Go 1.24 or higher

### Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/sigmanor/cherkasyoblenergo-api.git
   cd cherkasyoblenergo-api
   ```

2. **Configure environment variables**

   Create a `.env` file in the root directory. See the table below for all available options:

   | Variable | Required | Default | Description |
   |----------|----------|---------|-------------|
   | `DB_NAME` | No | `cherkasyoblenergo.db` | SQLite database file path |
   | `SERVER_PORT` | No | `8080` | Port for the API server |
   | `NEWS_URL` | No | `https://gita.cherkasyoblenergo.com/obl-main-controller/api/news2?size=18&category=1&page=0` | URL to parse schedules from |
   | `PARSING_INTERVAL_MINUTES` | No | `5` | How often to check for new schedules (minutes) |
   | `RATE_LIMIT_PER_MINUTE` | No | `60` | Max requests per minute per IP |
   | `CACHE_TTL_SECONDS` | No | `60` | Response cache duration in seconds |
   | `LOG_LEVEL` | No | `info` | Logging level (`debug`, `info`, `warn`, `error`) |
   | `FORCE_HTTPS` | No | `false` | Redirect HTTP to HTTPS |
   | `API_KEY` | No | - | If set, enables API key authentication |
   | `PROXY_MODE` | No | `none` | Proxy mode for real client IP detection: `cloudflare`, `standard`, or `none` |

3. **Run the application**

   ```bash
   go run ./cmd/server/main.go
   ```

### Building

```bash
go build -o cherkasyoblenergo_api ./cmd/server/main.go
```

## API Documentation

### Base URL

```
/cherkasyoblenergo/api
```

### Get Power Outage Schedule

```
GET /blackout-schedule
```

Retrieve scheduling records based on filter options.

#### Query Parameters

| Parameter | Required | Description |
|-----------|----------|-------------|
| `option` | Yes | `all`, `latest_n`, `by_date`, or `by_schedule_date` |
| `date` | For `by_date`, `by_schedule_date` | `YYYY-MM-DD`, `today`, or `tomorrow` |
| `limit` | For `latest_n` | Integer > 0 |
| `queue` | No | Comma-separated queue identifiers (e.g., `3_2` or `4_1,3_1`) |

#### Filter Options

- `all` - Retrieves all schedule records
- `latest_n` - Gets limited number of recent records (requires `limit`)
- `by_date` - Gets records for specific publication date (requires `date`)
- `by_schedule_date` - Gets records for specific schedule date (requires `date`, optional `limit`)

#### Example Requests

**Get all schedules:**
```bash
curl "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=all"
```

**Get latest 5 schedules:**
```bash
curl "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=5"
```

**Get today's schedule for queue 3_2:**
```bash
curl "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=by_schedule_date&date=today&queue=3_2"
```

**Get schedule with multiple queues:**
```bash
curl "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=1&queue=4_1,3_1,2_2"
```

#### Response Examples

**Full response (without queue filter):**
```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Schedule for November 14",
    "date": "2024-03-20T10:30:00Z",
    "schedule_date": "2024-11-14",
    "1_1": "08:00-10:00",
    "1_2": "10:00-12:00",
    "2_1": "12:00-14:00",
    "2_2": "14:00-16:00",
    "3_1": "09:00-11:00",
    "3_2": "11:00-13:00",
    "4_1": "13:00-15:00",
    "4_2": "15:00-17:00",
    "5_1": "07:00-09:00",
    "5_2": "09:00-11:00",
    "6_1": "11:00-13:00",
    "6_2": "13:00-15:00"
  }
]
```

**Filtered response (with queue filter):**
```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Schedule for November 14",
    "date": "2024-03-20T10:30:00Z",
    "schedule_date": "2024-11-14",
    "3_2": "00:30 - 02:30, 06:00 - 09:00"
  }
]
```

### Rate Limiting

- Default: 60 requests per minute per IP
- Configurable via `RATE_LIMIT_PER_MINUTE` environment variable
- Response headers:
  - `X-RateLimit-Limit` - Maximum requests allowed
  - `X-RateLimit-Remaining` - Remaining requests in current window

### Caching

- Responses are cached for 60 seconds by default
- Configurable via `CACHE_TTL_SECONDS` environment variable
- Response headers:
  - `Cache-Control` - Cache duration
  - `X-Cache` - `HIT` or `MISS`

### Authentication (Optional)

By default, the API is public and requires no authentication. For private instances, you can enable API key authentication:

1. Set the `API_KEY` environment variable:
   ```properties
   API_KEY=your-secret-key
   ```

2. Include the key in all requests:
   ```bash
   curl -H "X-API-Key: your-secret-key" "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=1"
   ```

If `API_KEY` is not set or empty, the API remains public.

### Error Responses

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Bad Request (invalid parameters) |
| 401 | Unauthorized (invalid or missing API key, if authentication is enabled) |
| 429 | Too Many Requests (rate limit exceeded) |
| 500 | Internal Server Error |

## Running Tests

```bash
go test ./...
```

## License

This project is licensed under the BSD 2-Clause License. See the [LICENSE](LICENSE) file for details.