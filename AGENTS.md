# AGENTS.md

This document provides guidelines for AI coding agents working on the cherkasyoblenergo-api project.

## Project Overview

A Go REST API that parses and serves electricity blackout schedules from Cherkasyoblenergo (Ukrainian energy provider). Built with Fiber v2 web framework, GORM ORM with SQLite, and scheduled news parsing via cron.

## Tech Stack

- **Language:** Go 1.24+
- **Web Framework:** Fiber v2
- **ORM:** GORM with SQLite driver
- **Testing:** Standard `testing` package + `testify` (assert/require)
- **Configuration:** Viper + godotenv
- **HTML Parsing:** goquery
- **Scheduling:** robfig/cron/v3

## Build & Run Commands

```bash
# Install dependencies
go mod download

# Build the application
go build -o cherkasyoblenergo_api ./cmd/server/main.go

# Build with version tag
go build -ldflags="-X 'cherkasyoblenergo-api/internal/config.AppVersion=v1.0.0'" -o cherkasyoblenergo_api ./cmd/server/main.go

# Run the application
go run ./cmd/server/main.go
```

## Testing Commands

```bash
# Run all tests
go test ./...

# Run all tests with verbose output
go test -v ./...

# Run a single test file
go test -v ./internal/handlers/schedule_test.go ./internal/handlers/schedule.go

# Run a single test function
go test -v -run TestGetSchedule_AllOption ./internal/handlers/

# Run tests for a specific package
go test -v ./internal/config/
go test -v ./internal/handlers/
go test -v ./internal/parser/
go test -v ./internal/utils/
go test -v ./internal/middleware/
go test -v ./internal/database/

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

## Project Structure

```
cherkasyoblenergo-api/
├── cmd/server/main.go          # Application entrypoint
├── internal/
│   ├── cache/                  # In-memory cache implementation
│   ├── config/                 # Configuration loading (Viper)
│   ├── database/               # Database connection setup
│   ├── handlers/               # HTTP request handlers
│   ├── logger/                 # Logging setup (slog)
│   ├── middleware/             # Fiber middleware (auth, rate limiting, logging)
│   ├── models/                 # GORM models
│   ├── parser/                 # HTML parsing and cron jobs
│   └── utils/                  # Utility functions (date extraction)
├── .env.example                # Environment variables template
└── docker-compose.app.yml      # Docker composition
```

## Code Style Guidelines

### Imports

Order imports in three groups separated by blank lines:
1. Standard library
2. External packages
3. Internal packages

```go
import (
    "context"
    "encoding/json"
    "time"

    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"

    "cherkasyoblenergo-api/internal/models"
    "cherkasyoblenergo-api/internal/utils"
)
```

### Naming Conventions

- **Packages:** lowercase, single word (e.g., `handlers`, `middleware`, `utils`)
- **Files:** snake_case (e.g., `schedule_cache.go`, `date_extractor.go`)
- **Test files:** `*_test.go` in same package
- **Structs:** PascalCase (e.g., `ScheduleCache`, `IPRateLimiter`)
- **Interfaces:** PascalCase, typically ending with `-er` suffix
- **Functions/Methods:** PascalCase for exported, camelCase for unexported
- **Variables:** camelCase (e.g., `scheduleCache`, `newsURL`)
- **Constants:** PascalCase or camelCase depending on export status

### Error Handling

- Return errors as the last return value
- Check errors immediately after function calls
- Use `fmt.Errorf` with `%w` for error wrapping
- Log errors with context before returning

```go
if err != nil {
    log.Printf("Failed to fetch data: %v", err)
    return fmt.Errorf("failed to initialize database: %w", err)
}
```

### Struct Tags

Use multiple tags on GORM models for database and JSON mapping:

```go
type Schedule struct {
    ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
    NewsID int    `gorm:"uniqueIndex" json:"news_id"`
    Title  string `gorm:"type:text" json:"title"`
}
```

### HTTP Handlers

Return `fiber.Handler` functions using closures for dependency injection:

```go
func GetSchedule(db *gorm.DB, cache *cache.ScheduleCache) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // handler logic
        return c.JSON(response)
    }
}
```

### Testing Patterns

- Use table-driven tests with descriptive names
- Use `testify/assert` and `testify/require` for assertions
- Create helper functions like `setupTestDB()` for test setup
- Use in-memory SQLite (`:memory:`) for database tests
- Test both success and error cases

```go
func TestFunction_ScenarioName(t *testing.T) {
    // Arrange
    db, cache := setupTestDB()
    
    // Act
    result, err := SomeFunction(input)
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expected, result)
}
```

### Middleware Pattern

Middleware functions return `fiber.Handler`:

```go
func SomeMiddleware(config Config) fiber.Handler {
    return func(c *fiber.Ctx) error {
        // pre-processing
        err := c.Next()
        // post-processing
        return err
    }
}
```

## Configuration

Environment variables (see `.env.example`):
- `DB_NAME` - SQLite database filename
- `SERVER_PORT` - HTTP server port (default: 8080)
- `NEWS_URL` - Source URL for schedule data
- `RATE_LIMIT_PER_MINUTE` - Rate limit (default: 60)
- `CACHE_TTL_SECONDS` - Cache TTL (default: 60)
- `LOG_LEVEL` - Logging level: debug, info, warn, error
- `API_KEY` - Optional API key for authentication
- `PROXY_MODE` - none, standard, or cloudflare

## Docker

```bash
# Build and run with Docker Compose
docker-compose -f docker-compose.app.yml up --build
```

## Important Notes

1. **Ukrainian text handling:** The codebase processes Ukrainian month names with fuzzy matching (Levenshtein distance) to handle typos
2. **Schedule queues:** Queues are numbered 1-6 with subqueues 1-2 (format: `X_Y`)
3. **Date formats:** API accepts `YYYY-MM-DD`, `today`, or `tomorrow`
4. **Caching:** Response caching is implemented at the handler level
5. **Rate limiting:** IP-based rate limiting stored in SQLite
