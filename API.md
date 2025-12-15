# API Documentation [![EN](https://img.shields.io/badge/English-0e7837.svg)](API.md) [![UA](https://img.shields.io/badge/Ukrainian-c9c9c9.svg)](API_UA.md)

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

- `option` (required): `all`, `latest_n`, `by_date`, or `by_schedule_date`
- `date` (required for `by_date` and `by_schedule_date`): `YYYY-MM-DD`, `today`, or `tomorrow`
- `limit` (required for `latest_n`, optional for `by_schedule_date`): Integer greater than 0
- `queue` (optional): Comma-separated queue identifiers (e.g., `3_2` or `4_1,3_1`)

#### Example Query (JSON equivalent shown for clarity)

```json
{
  "option": "all | latest_n | by_date | by_schedule_date",
  "date": "YYYY-MM-DD | today | tomorrow", // Required for "by_date" and "by_schedule_date" options
  "limit": 5, // Required for "latest_n", optional for "by_schedule_date"
  "queue": "3_2" // Optional: Filter by queue(s). Single: "3_2" or multiple: "4_1, 3_1" (comma-separated, X: 1-6, Y: 1-2)
}
```

#### Filter Options

- `all`: Retrieves all schedule records
- `latest_n`: Gets limited number of recent records (requires `limit`)
- `by_date`: Gets records for specific publication date (requires `date`)
- `by_schedule_date`: Gets records for specific schedule date extracted from title (requires `date`, optional `limit`). Results are sorted by publication date (most recent first).
- `queue` (optional): Filters response to include only the specified queue field(s). Accepts a single queue (e.g., "3_2") or multiple comma-separated queues (e.g., "4_1, 3_1"). Whitespace around commas is ignored. Duplicates are automatically removed.

#### Response

##### Full Response (without queue filter)

```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Schedule for November 14",
    "date": "2024-03-20T10:30:00Z",
    "schedule_date": "2024-11-14",
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
    "date": "2024-03-20T10:30:00Z",
    "schedule_date": "2024-11-14",
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
    "date": "2024-03-20T10:30:00Z",
    "schedule_date": "2024-11-14",
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
  "rate_limit": 6 // Optional: requests per minute (default: 6, must be > 0)
}
```

#### Response

```json
{
  "api_key": "ggj7d1slfkm",
  "rate_limit": 6
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

### Register Webhook

`POST /webhook`

Register a webhook URL to receive automatic notifications when new power outage schedules are parsed. The webhook will be triggered whenever new schedule data is available.

#### Request Body

```json
{
  "webhook_url": "https://your-domain.com/webhook-endpoint"
}
```

#### Response

```json
{
  "message": "Webhook registered successfully",
  "webhook_url": "https://your-domain.com/webhook-endpoint"
}
```

#### Webhook Payload Format

When a new schedule is available, your webhook will receive a POST request with the following payload. Each webhook notification contains exactly one schedule in the `schedules` array:

```json
{
  "schedules": [
    {
      "id": 1234,
      "news_id": 100,
      "title": "Schedule for November 14",
      "date": "2024-03-20T10:30:00Z",
      "schedule_date": "2024-11-14",
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
}
```

**Note**: The `schedules` array will always contain exactly one schedule object. If multiple new schedules are parsed simultaneously, your webhook endpoint will receive multiple separate POST requests, one for each schedule.

#### Webhook Headers

Your webhook endpoint will receive these headers:

- `Content-Type: application/json`
- `X-Token: YOUR_API_KEY`

### Delete Webhook

`DELETE /webhook`

Remove your webhook registration. This will disable automatic notifications for your API key.

#### Response

```json
{
  "message": "Webhook deleted successfully"
}
```

### Get Webhook Status

`GET /webhook`

Check the current status of your webhook registration.

#### Response

```json
{
  "webhook_url": "https://your-domain.com/webhook-endpoint",
  "webhook_enabled": true,
  "webhook_failed_attempts": 0
}
```

## Webhook Features

### Automatic Retry Logic

The webhook system includes automatic retry logic:
- **First attempt**: Immediate delivery
- **Second attempt**: 1 minute delay
- **Third attempt**: 10 minutes delay

### Automatic Disabling

If your webhook endpoint fails to respond successfully (HTTP 200) for 3 consecutive attempts, the webhook will be automatically disabled to prevent further delivery attempts.

### Integration with Schedule Parsing

Webhooks are automatically triggered whenever new power outage schedules are parsed from the source website. Each new schedule triggers a separate webhook notification immediately after being saved to the database, ensuring you receive real-time updates without polling the API. If multiple schedules are parsed in a single parsing cycle, your endpoint will receive multiple individual webhook calls, one per schedule.

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

### Get Schedules by Schedule Date

```bash
curl "https://api.example.com/cherkasyoblenergo/api/blackout-schedule?option=by_schedule_date&date=2025-12-05&limit=1&queue=4_1" \
  -H "X-API-Key: YOUR_API_KEY"
```

### Get Today's Schedule

```bash
curl "https://api.example.com/cherkasyoblenergo/api/blackout-schedule?option=by_schedule_date&date=today" \
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

### Register Webhook

```bash
curl -X POST "https://api.example.com/cherkasyoblenergo/api/webhook" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{"webhook_url":"https://your-domain.com/webhook-endpoint"}'
```

### Delete Webhook

```bash
curl -X DELETE "https://api.example.com/cherkasyoblenergo/api/webhook" \
  -H "X-API-Key: YOUR_API_KEY"
```

### Get Webhook Status

```bash
curl "https://api.example.com/cherkasyoblenergo/api/webhook" \
  -H "X-API-Key: YOUR_API_KEY"
```
