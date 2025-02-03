## Table of Contents

- [Table of Contents](#table-of-contents)
- [Description](#description)
- [Free access to API](#free-access-to-api)
- [How run](#how-run)
  - [Requirements](#requirements)
  - [Clone repo](#clone-repo)
  - [Config](#config)
  - [Deploy](#deploy)
- [API Documentation](#api-documentation)
  - [General Information](#general-information)
  - [Note](#note)
  - [Base URL](#base-url)
  - [Endpoints](#endpoints)
    - [/blackout-schedule](#blackout-schedule)
      - [Request](#request)
      - [Filter Options](#filter-options)
      - [Schedule JSON Response Object Fields](#schedule-json-response-object-fields)
      - [Example Request (latest\_n)](#example-request-latest_n)
      - [Response](#response)
    - [/generate-api-key](#generate-api-key)
      - [Request](#request-1)
      - [Query Parameters](#query-parameters)
      - [Example Request](#example-request)
      - [Response](#response-1)
    - [/update-api-key](#update-api-key)
      - [Request](#request-2)
      - [Query Parameters](#query-parameters-1)
      - [Example Requests](#example-requests)
      - [Response](#response-2)
- [Contribution](#contribution)
  - [run](#run)
  - [build](#build)

## Description

This unofficial API serves as a robust data gateway, delivering detailed power outage charts directly sourced from cherkasyoblenergo.com. It is designed to facilitate the seamless integration of real-time and historical power outage information into various applications and platforms.

## Free access to API
I am running this API on my own server for personal use. However, if you would like to test it or use it in your projects, I can grant you access without any uptime and others guarantees. The access will be limited to 2 requests per minute. If you're interested, please email me at sigmanor@pm.me

## How run

### Requirements

- [docker](https://docs.docker.com/engine/install/)
- [docker-compose](https://docs.docker.com/compose/install/)

### Clone repo

```shell
git clone https://github.com/Sigmanor/cherkasyoblenergo-api.git
```

### Config

Create a **.env** file in the app root directory:

```properties
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=your_strong_db_password
DB_NAME=myCoolDB
ADMIN_PASSWORD=your_strong_admin_password
SERVER_PORT=3000
```

### Deploy

-  `docker-compose --env-file .env up -d --build` or if you already has postgresql container running `docker-compose -f docker-compose.app-only.yml --env-file .env up -d --build`

## API Documentation

### General Information

- All responses are in JSON format.
- The API uses GORM for database operations, so database connectivity errors will result in HTTP 500 errors.
- The /generate-api-key endpoint is intended for administrative purposes and should be secured properly.
- Input validation and error handling are enforced to ensure correct parameters are passed.

### Note

Before using the API endpoints, ensure that the necessary middleware such as API key authentication, rate limiting, and logging are properly configured as per the server's setup in main.go.

### Base URL

```text
/cherkasyoblenergo/api
```

### Endpoints

#### /blackout-schedule

Retrieve scheduling records (blackout schedules) based on various filter options.

##### Request

```text
HTTP Method: POST
URL: /cherkasyoblenergo/api/blackout-schedule
Content-Type: application/json
```

Request Body JSON Structure:

```json
{
  "option": "all" // "latest_n" | "by_date",
  "date": "2022-02-24", // required only when option is "by_date"
  "limit": 5 // <integer> required only when option is "latest_n", must be > 0
}
```

##### Filter Options

- **option all**: Retrieves all schedule records.
- **option latest_n**: Retrieves a limited number of schedule records ordered by descending date. Request must include the **limit** parameter (integer > 0).
- **option by_date**: Retrieves schedule records for a specific date. Request must include the **date** parameter in format YYYY-MM-DD.

##### Schedule JSON Response Object Fields

- **id**: int64 (Unique schedule record identifier)
- **news_id**: int (Associated news ID)
- **title**: string (Schedule title)
- **date**: string (Date of the schedule in ISO 8601 format)
- **1_1**: string
- **1_2**: string
- **2_1**: string
- **2_2**: string
- **3_1**: string
- **3_2**: string
- **4_1**: string
- **4_2**: string
- **5_1**: string
- **5_2**: string
- **6_1**: string
- **6_2**: string

##### Example Request (latest_n)

```text
POST /cherkasyoblenergo/api/blackout-schedule
Content-Type: application/json

{
  "option": "latest_n",
  "limit": 5
}
```

##### Response

- **On Success:**
  - HTTP Status Code: `200 OK`
  - Response Body: JSON array of [objects](#schedule-json-response-object-fields)

- **On Failure (examples):**
  - If JSON is invalid:
    - HTTP Status Code: `400 Bad Request`
    - Response JSON: `{ "error": "Incorrect JSON format" }`
  - If "limit" is invalid for latest_n:
    - HTTP Status Code: `400 Bad Request`
    - Response JSON: `{ "error": "Invalid limit value, it must be greater than zero" }`
  - If "date" is missing or incorrectly formatted for by_date:
    - HTTP Status Code: `400 Bad Request`
    - Response JSON: `{ "error": "Date parameter is required in the format YYYY-MM-DD" }` or  
      `{ "error": "Invalid date format, expected YYYY-MM-DD" }`
  - For database errors:
    - HTTP Status Code: `500 Internal Server Error`
    - Appropriate error message in JSON.

#### /generate-api-key

Generates a new API key for accessing the API. This endpoint is protected and requires an admin password.

##### Request

```text
HTTP Method: GET
URL: /cherkasyoblenergo/api/generate-api-key
```

##### Query Parameters

- **admin_password (required)**: The administrative password to authorize key generation.
- **rate_limit (optional)**: The rate limit value for the key (defaults to 1 if not provided). It should be an integer.

##### Example Request

```text
GET /cherkasyoblenergo/api/generate-api-key?admin_password=<YOUR_ADMIN_PASSWORD>&rate_limit=5
```

##### Response

- **On Success:**
  - HTTP Status Code: `200 OK`
  - Response Body JSON Structure:
    ```json
    {
      "api_key": "ggj7d1slfkm",
      "rate_limit": 2
    }
    ```

- **On Failure:**
  - If admin_password does not match the configured value:
    - HTTP Status Code: `401 Unauthorized`
    - Response JSON: `{ "error": "Unauthorized" }`
  - If rate_limit is provided but invalid (non-integer):
    - HTTP Status Code: `400 Bad Request`
    - Response JSON: `{ "error": "Invalid rate_limit value" }`
  - If there are internal server errors (e.g., failure to create the API key in the database):
    - HTTP Status Code: `500 Internal Server Error`
    - Response JSON: `{ "error": "Failed to create API key" }`


#### /update-api-key

Allows an administrator to manage existing API keys

##### Request

```text
HTTP Method: GET
URL: /cherkasyoblenergo/api/update-api-key
```

##### Query Parameters

- **admin_password (required)**: The administrative password to authorize changes.
- **key (required)**: The existing API key to be managed.
- **update_key (optional, default: false)**: If set to `true`, generates and assigns a new API key.
- **delete_key (optional, default: false)**: If set to `true`, deletes the specified API key.
- **update_rate_limit (optional)**: Updates the rate limit for the specified API key. Should be an integer.

##### Example Requests

**Update API key:**
```text
GET /cherkasyoblenergo/api/update-api-key?admin_password=<YOUR_ADMIN_PASSWORD>&key=<EXISTING_KEY>&update_key=true
```

**Delete API key:**
```text
GET /cherkasyoblenergo/api/update-api-key?admin_password=<YOUR_ADMIN_PASSWORD>&key=<EXISTING_KEY>&delete_key=true
```

**Update rate limit:**
```text
GET /cherkasyoblenergo/api/update-api-key?admin_password=<YOUR_ADMIN_PASSWORD>&key=<EXISTING_KEY>&update_rate_limit=10
```

##### Response

- **On Success:**
  - HTTP Status Code: `200 OK`
  - Response Body (for key update):
    ```json
    {
      "message": "API key updated successfully",
      "new_key": "newly_generated_key"
    }
    ```
  - Response Body (for key deletion):
    ```json
    {
      "message": "API key deleted successfully"
    }
    ```
  - Response Body (for rate limit update):
    ```json
    {
      "message": "Rate limit updated successfully"
    }
    ```

- **On Failure:**
  - If admin_password is incorrect:
    - HTTP Status Code: `401 Unauthorized`
    - Response JSON: `{ "error": "Unauthorized" }`
  - If key is missing:
    - HTTP Status Code: `400 Bad Request`
    - Response JSON: `{ "error": "API key is required" }`
  - If key is not found:
    - HTTP Status Code: `404 Not Found`
    - Response JSON: `{ "error": "API key not found" }`
  - If update_rate_limit is invalid:
    - HTTP Status Code: `400 Bad Request`
    - Response JSON: `{ "error": "Invalid rate_limit value" }`
  - If internal server error occurs:
    - HTTP Status Code: `500 Internal Server Error`
    - Response JSON: `{ "error": "Failed to update API key" }`



## Contribution

- Install [go](https://go.dev/dl/)
- Install [postgresql](https://www.postgresql.org/)

### run

```shell
go run ./cmd/server/main.go
```

### build

```shell
go build -o cherkasyoblenergo_api ./cmd/server/main.go
``` 