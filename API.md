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

`GET /blackout-schedule`

Retrieve scheduling records based on filter options.

#### Query Parameters

- `option` (required): `all`, `latest_n`, or `by_date`
- `date` (required for `by_date`): `YYYY-MM-DD`
- `limit` (required for `latest_n`): Integer greater than 0
- `queue` (optional): Comma-separated queue identifiers (e.g., `3_2` or `4_1,3_1`)

#### Example Query (JSON equivalent shown for clarity)

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

### Create API Key

`POST /api-keys`

Generate a new API key with optional rate limiting. Provide the admin password in the request body.

#### Request Body

```json
{
  "admin_password": "YOUR_ADMIN_PASSWORD",
  "rate_limit": 2 // Optional: requests per minute (default: 2, must be > 0)
}
```

#### Response

```json
{
  "api_key": "ggj7d1slfkm",
  "rate_limit": 2
}
```

### Update API Key

`PATCH /api-keys`

Rotate an existing API key and/or update its rate limit. Provide credentials and the target key in the request body.

#### Request Body

```json
{
  "admin_password": "YOUR_ADMIN_PASSWORD",
  "key": "target_api_key",
  "rotate_key": true, // Optional: generate a new key value
  "rate_limit": 5 // Optional: new rate limit, must be > 0
}
```

At least one field is required.

#### Response

```json
{
  "message": "API key updated successfully",
  "new_key": "new_generated_key",
  "new_rate_limit": 5
}
```

### Delete API Key

`DELETE /api-keys`

Remove an API key permanently. Supply the admin password and target key in the request body.

#### Request Body

```json
{
  "admin_password": "YOUR_ADMIN_PASSWORD",
  "key": "target_api_key"
}
```

#### Response

```json
{
  "message": "API key deleted successfully"
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

### Get Latest Schedules (GET)

```bash
curl "https://api.example.com/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=5" \
  -H "X-API-Key: YOUR_API_KEY"
```

### Get Latest Schedules with Queue Filter

```bash
curl "https://api.example.com/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=5&queue=3_2" \
  -H "X-API-Key: YOUR_API_KEY"
```

### Get Latest Schedules with Multiple Queue Filter

```bash
curl "https://api.example.com/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=5&queue=4_1,3_1,2_2" \
  -H "X-API-Key: YOUR_API_KEY"
```

### Create API Key

```bash
curl -X POST "https://api.example.com/cherkasyoblenergo/api/api-keys" \
  -H "Content-Type: application/json" \
  -d '{"admin_password":"YOUR_ADMIN_PASSWORD","rate_limit":5}'
```

### Update API Key

```bash
curl -X PATCH "https://api.example.com/cherkasyoblenergo/api/api-keys" \
  -H "Content-Type: application/json" \
  -d '{"admin_password":"YOUR_ADMIN_PASSWORD","key":"YOUR_KEY","rotate_key":true,"rate_limit":4}'
```

### Delete API Key

```bash
curl -X DELETE "https://api.example.com/cherkasyoblenergo/api/api-keys" \
  -H "Content-Type: application/json" \
  -d '{"admin_password":"YOUR_ADMIN_PASSWORD","key":"YOUR_KEY"}'
```
