# Cherkasyoblenergo API Documentation

## Authentication

All API requests require an API key to be included in the header:
```
Authorization: Bearer YOUR_API_KEY
```

## Rate Limits

- Free tier: 2 requests/minute
- Premium tier: Contact for custom limits

## Endpoints

### Get Power Outage Schedule
`POST /blackout-schedule`

Request body:
```json
{
  "city": "string",
  "street": "string",
  "date": "YYYY-MM-DD"
}
```

Response:
```json
{
  "schedule": [
    {
      "startTime": "HH:MM",
      "endTime": "HH:MM",
      "status": "active|planned"
    }
  ]
}
```

### Generate API Key (Admin)
`GET /generate-api-key`

Headers:
```
Authorization: Bearer ADMIN_KEY
```

Response:
```json
{
  "apiKey": "string",
  "expiresAt": "YYYY-MM-DD"
}
```

### Update API Key (Admin)
`GET /update-api-key`

Headers:
```
Authorization: Bearer ADMIN_KEY
```

Query parameters:
- `key` - API key to update
- `action` - Action to perform (revoke|extend)

Response:
```json
{
  "status": "success|error",
  "message": "string"
}
```

## Error Codes

- 401: Unauthorized - Invalid or missing API key
- 403: Forbidden - Rate limit exceeded
- 422: Validation Error - Invalid request parameters
- 500: Internal Server Error

## Examples

### Get Schedule Example

```bash
curl -X POST https://api.example.com/cherkasyoblenergo/api/blackout-schedule \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "city": "Cherkasy",
    "street": "Shevchenko",
    "date": "2024-03-20"
  }'
```