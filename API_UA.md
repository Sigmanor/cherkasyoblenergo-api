# Документація Cherkasyoblenergo API [![EN](https://img.shields.io/badge/English-c9c9c9.svg)](API.md) [![UA](https://img.shields.io/badge/Ukrainian-0e7837.svg)](API_UA.md)

## Базовий URL

```
/cherkasyoblenergo/api
```

## Автентифікація

Всі ендпоінти вимагають API ключ. Включіть його у ваші запити як заголовок:

```
X-API-Key: YOUR_API_KEY
```

## Ендпоінти

### Отримати графік відключень електроенергії

`POST /blackout-schedule`

Отримати записи розкладу на основі параметрів фільтрації.

#### Тіло запиту

```json
{
  "option": "all | latest_n | by_date",
  "date": "YYYY-MM-DD", // Обов'язково для опції "by_date"
  "limit": 5 // Обов'язково для опції "latest_n", має бути > 0
}
```

#### Опції фільтрації

- `all`: Отримує всі записи розкладу
- `latest_n`: Отримує обмежену кількість останніх записів (потрібен `limit`)
- `by_date`: Отримує записи для конкретної дати (потрібна `date`)

#### Відповідь

```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Назва розкладу",
    "date": "2024-03-20",
    "1_1": "текст",
    "1_2": "текст",
    "2_1": "текст",
    "2_2": "текст",
    "3_1": "текст",
    "3_2": "текст",
    "4_1": "текст",
    "4_2": "текст",
    "5_1": "текст",
    "5_2": "текст",
    "6_1": "текст",
    "6_2": "текст"
  }
]
```

### Згенерувати API ключ

`GET /generate-api-key`

Згенерувати новий API ключ з опціональним обмеженням швидкості.

#### Параметри запиту

- `admin_password` (обов'язково): Адміністративний пароль
- `rate_limit` (опціонально): Запитів на хвилину (за замовчуванням: 1)

#### Відповідь

```json
{
  "api_key": "ggj7d1slfkm",
  "rate_limit": 2
}
```

### Оновити API ключ

`GET /update-api-key`

Керувати існуючими API ключами.

#### Параметри запиту

- `admin_password` (обов'язково): Адміністративний пароль
- `key` (обов'язково): API ключ для керування
- `update_key` (опціонально): Встановіть true для генерації нового ключа
- `delete_key` (опціонально): Встановіть true для видалення ключа
- `update_rate_limit` (опціонально): Нове значення обмеження швидкості

#### Приклади відповідей

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

## Обробка помилок

| Status Code | Description           | Example Response                         |
| ----------- | --------------------- | ---------------------------------------- |
| 200         | Success               | `{"message": "Success"}`                 |
| 400         | Bad Request           | `{"error": "Invalid JSON format"}`       |
| 401         | Unauthorized          | `{"error": "Unauthorized"}`              |
| 404         | Not Found             | `{"error": "API key not found"}`         |
| 500         | Internal Server Error | `{"error": "Failed to process request"}` |

## Приклади використання

### Отримати останні розклади

```bash
curl -X POST "https://api.example.com/cherkasyoblenergo/api/blackout-schedule" \
  -H "X-API-Key: YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "option": "latest_n",
    "limit": 5
  }'
```

### Згенерувати API ключ

```bash
curl "https://api.example.com/cherkasyoblenergo/api/generate-api-key?admin_password=YOUR_ADMIN_PASSWORD&rate_limit=5"
```
