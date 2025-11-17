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

`GET /blackout-schedule`

Отримати записи розкладу на основі параметрів фільтрації.

#### Параметри запиту

- `option` (обов'язково): `all`, `latest_n` або `by_date`
- `date` (обов'язково для `by_date`): `YYYY-MM-DD`
- `limit` (обов'язково для `latest_n`): ціле число > 0
- `queue` (опціонально): значення черг через кому (наприклад, `3_2` або `4_1,3_1`)

#### Приклад параметрів (JSON еквівалент для наочності)

```json
{
  "option": "all | latest_n | by_date",
  "date": "YYYY-MM-DD", // Обов'язково для опції "by_date"
  "limit": 5, // Обов'язково для опції "latest_n", має бути > 0
  "queue": "3_2" // Опціонально: Фільтр по черзі(чергах). Одна: "3_2" або декілька: "4_1, 3_1" (через кому, X: 1-6, Y: 1-2)
}
```

#### Опції фільтрації

- `all`: Отримує всі записи розкладу
- `latest_n`: Отримує обмежену кількість останніх записів (потрібен `limit`)
- `by_date`: Отримує записи для конкретної дати (потрібна `date`)
- `queue` (опціонально): Фільтрує відповідь, щоб включити лише вказане поле(поля) черги. Приймає одну чергу (наприклад, "3_2") або декілька черг через кому (наприклад, "4_1, 3_1"). Пробіли навколо ком ігноруються. Дублікати автоматично видаляються.

#### Відповідь

#### Приклади відповідей

**Повна відповідь (без фільтра черги):**

```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Графік на 14 листопада",
    "date": "2024-03-20",
    "schedule_date": "14.11",
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

**Відфільтрована відповідь (з фільтром черги):**

```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Графік на 14 листопада",
    "date": "2024-03-20",
    "schedule_date": "14.11",
    "3_2": "00:30 - 02:30, 06:00 - 09:00"
  }
]
```

**Відфільтрована відповідь (з декількома чергами):**

```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Графік на 14 листопада",
    "date": "2024-03-20",
    "schedule_date": "14.11",
    "4_1": "01:00 - 03:00, 06:00 - 09:00",
    "3_1": "00:30 - 02:30, 05:30 - 08:00"
  }
]
```

### Створити API ключ

`POST /api-keys`

Створити новий API ключ з опціональним обмеженням швидкості. Передайте `admin_password` у тілі запиту.

#### Тіло запиту

```json
{
  "admin_password": "YOUR_ADMIN_PASSWORD",
  "rate_limit": 6 // Опціонально: запитів на хвилину (за замовчуванням 6, має бути > 0)
}
```

#### Відповідь

```json
{
  "api_key": "ggj7d1slfkm",
  "rate_limit": 6
}
```

### Оновити API ключ

`PATCH /api-keys`

Повернути новий ключ або змінити ліміт швидкості. Передайте `admin_password` і цільовий ключ у тілі запиту.

#### Тіло запиту

```json
{
  "admin_password": "YOUR_ADMIN_PASSWORD",
  "key": "target_api_key",
  "rotate_key": true, // Опціонально: згенерувати нове значення ключа
  "rate_limit": 5 // Опціонально: новий ліміт запитів (> 0)
}
```

Потрібно вказати хоча б одне поле.

#### Відповідь

```json
{
  "message": "API key updated successfully",
  "new_key": "new_generated_key",
  "new_rate_limit": 5
}
```

### Видалити API ключ

`DELETE /api-keys`

Повністю видалити API ключ. Передайте `admin_password` і ключ у тілі запиту.

#### Тіло запиту

```json
{
  "admin_password": "YOUR_ADMIN_PASSWORD",
  "key": "target_api_key"
}
```

#### Відповідь

```json
{
  "message": "API key deleted successfully"
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
curl "https://api.example.com/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=5" \
  -H "X-API-Key: YOUR_API_KEY"
```

### Отримати останні розклади з фільтром черги

```bash
curl "https://api.example.com/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=5&queue=3_2" \
  -H "X-API-Key: YOUR_API_KEY"
```

### Отримати останні розклади з фільтром декількох черг

```bash
curl "https://api.example.com/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=5&queue=4_1,3_1,2_2" \
  -H "X-API-Key: YOUR_API_KEY"
```

### Створити API ключ

```bash
curl -X POST "https://api.example.com/cherkasyoblenergo/api/api-keys" \
  -H "Content-Type: application/json" \
  -d '{"admin_password":"YOUR_ADMIN_PASSWORD","rate_limit":5}'
```

### Оновити API ключ

```bash
curl -X PATCH "https://api.example.com/cherkasyoblenergo/api/api-keys" \
  -H "Content-Type: application/json" \
  -d '{"admin_password":"YOUR_ADMIN_PASSWORD","key":"YOUR_KEY","rotate_key":true,"rate_limit":4}'
```

### Видалити API ключ

```bash
curl -X DELETE "https://api.example.com/cherkasyoblenergo/api/api-keys" \
  -H "Content-Type: application/json" \
  -d '{"admin_password":"YOUR_ADMIN_PASSWORD","key":"YOUR_KEY"}'
```
