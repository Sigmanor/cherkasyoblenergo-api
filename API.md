# Cherkasyoblenergo API Documentation

## Base URL
```
/cherkasyoblenergo/api
```

## Authentication
All endpoints require an API key. Include it in your requests as a header:
```
Authorization: Bearer YOUR_API_KEY
```

## Endpoints

### Get Power Outage Schedule
`POST /blackout-schedule`

Retrieve scheduling records based on filter options.

#### Request Body
```json
{
  "option": "all | latest_n | by_date",
  "date": "YYYY-MM-DD",     // Required for "by_date" option
  "limit": 5                // Required for "latest_n" option, must be > 0
}
```

#### Filter Options
- `all`: Retrieves all schedule records
- `latest_n`: Gets limited number of recent records (requires `limit`)
- `by_date`: Gets records for specific date (requires `date`)

#### Response
```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Schedule Title",
    "date": "2024-03-20",
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

| Status Code | Description            | Example Response                                        |
|-------------|------------------------|--------------------------------------------------------|
| 200         | Success               | `{"message": "Success"}`                                |
| 400         | Bad Request           | `{"error": "Invalid JSON format"}`                      |
| 401         | Unauthorized          | `{"error": "Unauthorized"}`                             |
| 404         | Not Found             | `{"error": "API key not found"}`                        |
| 500         | Internal Server Error | `{"error": "Failed to process request"}`                |

## Example Usage

### Get Latest Schedules
```bash
curl -X POST "https://api.example.com/cherkasyoblenergo/api/blackout-schedule" \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "option": "latest_n",
    "limit": 5
  }'
```

### Generate API Key
```bash
curl "https://api.example.com/cherkasyoblenergo/api/generate-api-key?admin_password=YOUR_ADMIN_PASSWORD&rate_limit=5"
```