# Cherkasyoblenergo API Documentation [![EN](https://img.shields.io/badge/English-0e7837.svg)](API.md) [![UA](https://img.shields.io/badge/Ukrainian-c9c9c9.svg)](API_UA.md)

## Base URL

```
/cherkasyoblenergo/api
```

## Authentication

All endpoints require an API key. Include it in your requests as a header:

```
X-API-Key: YOUR_API_KEY
```

## Endpoints

### Get Power Outage Schedule

`POST /blackout-schedule`

Retrieve scheduling records based on filter options.

#### Request Body

```json
{
  "option": "all | latest_n | by_date",
  "date": "YYYY-MM-DD", // Required for "by_date" option
  "limit": 5, // Required for "latest_n" option, must be > 0
  "queue": "3_2" // Optional: Filter by queue(s). Single: "3_2" or multiple: "4_1, 3_1" (comma-separated, X: 1-6, Y: 1-2)
}
```

#### Filter Options

- `all`: Retrieves all schedule records
- `latest_n`: Gets limited number of recent records (requires `limit`)
- `by_date`: Gets records for specific date (requires `date`)
- `queue` (optional): Filters response to include only the specified queue field(s). Accepts a single queue (e.g., "3_2") or multiple comma-separated queues (e.g., "4_1, 3_1"). Whitespace around commas is ignored. Duplicates are automatically removed.

#### Response

##### Full Response (without queue filter)

```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Schedule for November 14",
    "date": "2024-03-20",
    "schedule_date": "14.11",
    "1_1": "text",
    "1_2": "text",
    "2_1": "text",
    "2_2": "text",
    "3_1": "text",
    "3_2": "text",
    "4_1": "text",
    "4_2": "text",
    "5_1": "text",
    "5_2": "text",
    "6_1": "text",
    "6_2": "text"
  }
]
```

##### Filtered Response (with queue filter)

```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Schedule for November 14",
    "date": "2024-03-20",
    "schedule_date": "14.11",
    "3_2": "00:30 - 02:30, 06:00 - 09:00"
  }
]
```

##### Filtered Response (with multiple queues)

```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Schedule for November 14",
    "date": "2024-03-20",
    "schedule_date": "14.11",
    "4_1": "01:00 - 03:00, 06:00 - 09:00",
    "3_1": "00:30 - 02:30, 05:30 - 08:00"
  }
]
```

### Generate API Key

`GET /generate-api-key`

Generate new API key with optional rate limiting.

#### Query Parameters

- `admin_password` (required): Administrative password
- `rate_limit` (optional): Requests per minute (default: 1)

#### Response

```json
{
  "api_key": "ggj7d1slfkm",
  "rate_limit": 2
}
```

### Update API Key

`GET /update-api-key`

Manage existing API keys.

#### Query Parameters

- `admin_password` (required): Administrative password
- `key` (required): API key to manage
- `update_key` (optional): Set true to generate new key
- `delete_key` (optional): Set true to delete key
- `update_rate_limit` (optional): New rate limit value

#### Response Examples

```json
// Update key
{
  "message": "API key updated successfully",
  "new_key": "new_generated_key"
}

// Delete key
{
  "message": "API key deleted successfully"
}

// Update rate limit
{
  "message": "Rate limit updated successfully"
}
```

## Error Handling

| Status Code | Description           | Example Response                         |
| ----------- | --------------------- | ---------------------------------------- |
| 200         | Success               | `{"message": "Success"}`                 |
| 400         | Bad Request           | `{"error": "Invalid JSON format"}`       |
| 401         | Unauthorized          | `{"error": "Unauthorized"}`              |
| 404         | Not Found             | `{"error": "API key not found"}`         |
| 500         | Internal Server Error | `{"error": "Failed to process request"}` |

## Example Usage

### Get Latest Schedules

```bash
curl -X POST "https://api.example.com/cherkasyoblenergo/api/blackout-schedule" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "option": "latest_n",
    "limit": 5
  }'
```

### Get Latest Schedules with Queue Filter

```bash
curl -X POST "https://api.example.com/cherkasyoblenergo/api/blackout-schedule" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "option": "latest_n",
    "limit": 5,
    "queue": "3_2"
  }'
```

### Get Latest Schedules with Multiple Queue Filter

```bash
curl -X POST "https://api.example.com/cherkasyoblenergo/api/blackout-schedule" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "option": "latest_n",
    "limit": 5,
    "queue": "4_1, 3_1, 2_2"
  }'
```

### Generate API Key

```bash
curl "https://api.example.com/cherkasyoblenergo/api/generate-api-key?admin_password=YOUR_ADMIN_PASSWORD&rate_limit=5"
```
