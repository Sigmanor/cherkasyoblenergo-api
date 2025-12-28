<div align="center">

# Cherkasyoblenergo API

[![EN](https://img.shields.io/badge/English-c9c9c9.svg)](README.md) [![UA](https://img.shields.io/badge/Ukrainian-0e7837.svg)](README_UA.md)

[![Deploy](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml)
[![Tests](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Sigmanor/cherkasyoblenergo-api)](https://go.dev/)
[![License](https://img.shields.io/github/license/Think-Root/chappie_server)](LICENSE)
[![Releases](https://img.shields.io/github/release/Sigmanor/cherkasyoblenergo-api.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/releases)
[![Changelog](https://img.shields.io/badge/changelog-md-blue)](CHANGELOG.md)

**Відкритий публічний API** для отримання графіків відключень електроенергії з [cherkasyoblenergo.com](https://cherkasyoblenergo.com/)

</div>

## Публічний інстанс

Публічний інстанс цього API доступний за адресою:

```
https://hue.pp.ua/cherkasyoblenergo/api/
```

Ви можете вільно використовувати його для своїх проєктів. Однак зверніть увагу, що цей інстанс надається "як є" без будь-яких гарантій доступності чи безперебійної роботи. Для продакшн використання або якщо вам потрібна надійність, рекомендуємо розгорнути власний інстанс.

## Можливості

- Дані графіків відключень електроенергії в реальному часі
- Доступ до історичних даних
- RESTful API інтерфейс
- Rate limiting по IP адресі
- Кешування відповідей
- Опціональна автентифікація за API ключем

## Встановлення

### Передумови

- Go 1.24 або вище
- PostgreSQL 17

### Налаштування

1. **Встановіть PostgreSQL 17**

   Дотримуйтесь [офіційного посібника з встановлення PostgreSQL](https://www.postgresql.org/download/).

2. **Клонуйте репозиторій**

   ```bash
   git clone https://github.com/Sigmanor/cherkasyoblenergo-api.git
   cd cherkasyoblenergo-api
   ```

3. **Налаштуйте змінні середовища**

   Створіть файл `.env` у кореневій директорії. Дивіться таблицю нижче для всіх доступних опцій:

   | Змінна | Обов'язкова | За замовчуванням | Опис |
   |--------|-------------|------------------|------|
   | `DB_HOST` | Так | - | Хост PostgreSQL |
   | `DB_PORT` | Так | - | Порт PostgreSQL |
   | `DB_USER` | Так | - | Ім'я користувача PostgreSQL |
   | `DB_PASSWORD` | Так | - | Пароль PostgreSQL |
   | `DB_NAME` | Так | - | Назва бази даних PostgreSQL |
   | `SERVER_PORT` | Ні | `8080` | Порт для API сервера |
   | `NEWS_URL` | Ні | `https://gita.cherkasyoblenergo.com/obl-main-controller/api/news2?size=18&category=1&page=0` | URL для парсингу графіків |
   | `PARSING_INTERVAL_MINUTES` | Ні | `5` | Як часто перевіряти нові графіки (хвилини) |
   | `RATE_LIMIT_PER_MINUTE` | Ні | `60` | Макс. запитів на хвилину на IP |
   | `CACHE_TTL_SECONDS` | Ні | `60` | Тривалість кешу відповідей в секундах |
   | `LOG_LEVEL` | Ні | `info` | Рівень логування (`debug`, `info`, `warn`, `error`) |
   | `FORCE_HTTPS` | Ні | `false` | Перенаправляти HTTP на HTTPS |
   | `API_KEY` | Ні | - | Якщо встановлено, вмикає автентифікацію за API ключем |
   | `PROXY_MODE` | Ні | `none` | Режим проксі для визначення реального IP клієнта: `cloudflare`, `standard` або `none` |

4. **Запустіть додаток**

   ```bash
   go run ./cmd/server/main.go
   ```

### Збірка

```bash
go build -o cherkasyoblenergo_api ./cmd/server/main.go
```

## Документація API

### Базовий URL

```
/cherkasyoblenergo/api
```

### Отримати графік відключень

```
GET /blackout-schedule
```

Отримати записи розкладу на основі параметрів фільтрації.

#### Параметри запиту

| Параметр | Обов'язковий | Опис |
|----------|--------------|------|
| `option` | Так | `all`, `latest_n`, `by_date` або `by_schedule_date` |
| `date` | Для `by_date`, `by_schedule_date` | `YYYY-MM-DD`, `today` або `tomorrow` |
| `limit` | Для `latest_n` | Ціле число > 0 |
| `queue` | Ні | Значення черг через кому (наприклад, `3_2` або `4_1,3_1`) |

#### Опції фільтрації

- `all` - Отримує всі записи розкладу
- `latest_n` - Отримує обмежену кількість останніх записів (потрібен `limit`)
- `by_date` - Отримує записи за датою публікації (потрібна `date`)
- `by_schedule_date` - Отримує записи за датою графіка (потрібна `date`, опціонально `limit`)

#### Приклади запитів

**Отримати всі графіки:**
```bash
curl "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=all"
```

**Отримати останні 5 графіків:**
```bash
curl "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=5"
```

**Отримати сьогоднішній графік для черги 3_2:**
```bash
curl "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=by_schedule_date&date=today&queue=3_2"
```

**Отримати графік з кількома чергами:**
```bash
curl "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=1&queue=4_1,3_1,2_2"
```

#### Приклади відповідей

**Повна відповідь (без фільтра черги):**
```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Графік на 14 листопада",
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

**Відфільтрована відповідь (з фільтром черги):**
```json
[
  {
    "id": 1234,
    "news_id": 100,
    "title": "Графік на 14 листопада",
    "date": "2024-03-20T10:30:00Z",
    "schedule_date": "2024-11-14",
    "3_2": "00:30 - 02:30, 06:00 - 09:00"
  }
]
```

### Rate Limiting

- За замовчуванням: 60 запитів на хвилину на IP
- Налаштовується через змінну середовища `RATE_LIMIT_PER_MINUTE`
- Заголовки відповіді:
  - `X-RateLimit-Limit` - Максимальна кількість запитів
  - `X-RateLimit-Remaining` - Залишок запитів у поточному вікні

### Кешування

- Відповіді кешуються на 60 секунд за замовчуванням
- Налаштовується через змінну середовища `CACHE_TTL_SECONDS`
- Заголовки відповіді:
  - `Cache-Control` - Тривалість кешу
  - `X-Cache` - `HIT` або `MISS`

### Автентифікація (опціонально)

За замовчуванням API є публічним і не вимагає автентифікації. Для приватних інстансів можна увімкнути автентифікацію за API ключем:

1. Встановіть змінну середовища `API_KEY`:
   ```properties
   API_KEY=your-secret-key
   ```

2. Додавайте ключ до всіх запитів:
   ```bash
   curl -H "X-API-Key: your-secret-key" "http://localhost:8080/cherkasyoblenergo/api/blackout-schedule?option=latest_n&limit=1"
   ```

Якщо `API_KEY` не встановлено або порожнє, API залишається публічним.

### Коди помилок

| Код статусу | Опис |
|-------------|------|
| 200 | Успіх |
| 400 | Невірний запит (невалідні параметри) |
| 401 | Неавторизовано (невірний або відсутній API ключ, якщо автентифікацію увімкнено) |
| 429 | Забагато запитів (перевищено rate limit) |
| 500 | Внутрішня помилка сервера |

## Запуск тестів

```bash
go test ./...
```

## Ліцензія

Цей проєкт ліцензовано BSD 2-Clause. Деталі дивіться у файлі [LICENSE](LICENSE).