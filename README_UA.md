<div align="center">

# Cherkasyoblenergo API

[![EN](https://img.shields.io/badge/English-c9c9c9.svg)](README.md) [![UA](https://img.shields.io/badge/Ukrainian-0e7837.svg)](README_UA.md)

[![Deploy](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/deploy.yml)
[![Tests](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml/badge.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/actions/workflows/tests.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/Sigmanor/cherkasyoblenergo-api)](https://go.dev/)
[![License](https://img.shields.io/github/license/Think-Root/chappie_server)](LICENSE)
[![Releases](https://img.shields.io/github/release/Sigmanor/cherkasyoblenergo-api.svg)](https://github.com/Sigmanor/cherkasyoblenergo-api/releases)
[![Changelog](https://img.shields.io/badge/changelog-md-blue)](CHANGELOG.md)

Неофіційний API сервіс для отримання графіків відключень електроенергії з [cherkasyoblenergo.com](https://cherkasyoblenergo.com/). Отримуйте інформацію про відключення електроенергії в реальному часі та історичні дані через RESTful API інтерфейс.

</div>

## 📋 Зміст

- [Cherkasyoblenergo API](#cherkasyoblenergo-api)
  - [📋 Зміст](#-зміст)
  - [✨ Ключові особливості](#-ключові-особливості)
  - [🚀 Встановлення](#-встановлення)
    - [Передумови](#передумови)
    - [Налаштування](#налаштування)
  - [🔑 Документація API](#-документація-api)
    - [Базовий URL](#базовий-url)
    - [Доступні ендпоінти](#доступні-ендпоінти)
  - [💻 Розробка](#-розробка)
    - [Вимоги](#вимоги)
    - [Локальна розробка](#локальна-розробка)
  - [🤝 Внесок у проект](#-внесок-у-проект)
  - [❗ Усунення неполадок](#-усунення-неполадок)
  - [🚦 Запуск тестів](#-запуск-тестів)
  - [⚡ Безкоштовний доступ до API](#-безкоштовний-доступ-до-api)

## ✨ Ключові особливості

- Дані графіків відключень електроенергії в реальному часі
- Доступ до історичних даних
- RESTful API інтерфейс
- Підтримка обмеження швидкості запитів
- Автентифікація за допомогою API ключа
- Підтримка розгортання через Docker

## 🚀 Встановлення

### Передумови

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- PostgreSQL 17 (тільки якщо запускаєте без Docker)

### Налаштування

1. Клонуйте репозиторій:

```bash
git clone https://github.com/Sigmanor/cherkasyoblenergo-api.git
cd cherkasyoblenergo-api
```

2. Створіть файл `.env` з необхідними конфігураціями:

```properties
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=your_strong_db_password
DB_NAME=myCoolDB
ADMIN_PASSWORD=your_strong_admin_password
SERVER_PORT=3000
```

3. Оберіть метод розгортання:

**Повне розгортання через Docker (з PostgreSQL):**

```bash
# Створіть постійний том для PostgreSQL
docker volume create postgres_data

# Розгорніть додаток і базу даних
docker compose -f docker-compose.app.yml -f docker-compose.db.yml up -d --build
```

**Розгортання тільки додатку (для існуючої PostgreSQL):**

```bash
docker compose -f docker-compose.app.yml up -d --build
```

## 🔑 Документація API

### Базовий URL

```
/cherkasyoblenergo/api
```

### Доступні ендпоінти

- `POST /blackout-schedule` - Отримати графіки відключень електроенергії
- `GET /generate-api-key` - Згенерувати API ключ (тільки для адміністратора)
- `GET /update-api-key` - Керувати API ключами (тільки для адміністратора)

[Детальна документація API](API_UA.md)

## 💻 Розробка

### Вимоги

- Go 1.23 або вище
- PostgreSQL 17
- Docker та Docker Compose (для контейнеризованого розгортання)

### Локальна розробка

```bash
# Запуск локально
go run ./cmd/server/main.go

# Збірка
go build -o cherkasyoblenergo_api ./cmd/server/main.go
```

## 🤝 Внесок у проект

1. Зробіть форк репозиторію
2. Створіть гілку для вашої функції (`git checkout -b feature/amazing-feature`)
3. Зафіксуйте ваші зміни (`git commit -m 'Add amazing feature'`)
4. Відправте зміни в гілку (`git push origin feature/amazing-feature`)
5. Відкрийте Pull Request

## ❗ Усунення неполадок

- **Проблеми з підключенням до бази даних**:
  - Для розгортання через Docker: Перевірте, чи створено том postgres_data
  - Перевірте облікові дані PostgreSQL та налаштування підключення в `.env`
  - Для повного налаштування Docker переконайтеся, що сервіс db працює правильно
- **Проблеми з API ключем**: Переконайтеся в правильній генерації API ключа та конфігурації обмеження швидкості
- **Проблеми з Docker**:
  - Перевірте логи Docker: `docker-compose logs`
  - Перевірте конфігурацію мережі Docker
  - Переконайтеся, що всі необхідні змінні середовища встановлені

## 🚦 Запуск тестів

Для запуску тестів локально:

```bash
go test ./...
```

## ⚡ Безкоштовний доступ до API

Оскільки цей API сервер постійно запущений у мене на хостингу (для моїх потреб) я можу надати вам обмежений доступ (2 запити/хв) до нього безкоштовно. Зв'яжіться зі мною за допомогою [email](mailto:dock-brunt-rarity@duck.com) для отримання доступу.
