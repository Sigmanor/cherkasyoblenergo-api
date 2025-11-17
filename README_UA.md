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
    - [Збірка](#збірка)
  - [🔑 Документація API](#-документація-api)
    - [Базовий URL](#базовий-url)
    - [Доступні ендпоінти](#доступні-ендпоінти)
  - [🚦 Запуск тестів](#-запуск-тестів)
  - [⚡ Безкоштовний доступ до API](#-безкоштовний-доступ-до-api)

## ✨ Ключові особливості

- Дані графіків відключень електроенергії в реальному часі
- Доступ до історичних даних
- RESTful API інтерфейс
- Підтримка обмеження швидкості запитів
- Автентифікація за допомогою API ключа

## 🚀 Встановлення

### Передумови

- Go 1.23 або вище
- PostgreSQL 17

### Налаштування

1. **Встановіть PostgreSQL 17**

   Дотримуйтесь [офіційного посібника з встановлення PostgreSQL](https://www.postgresql.org/download/), щоб встановити PostgreSQL на вашій системі.

2. **Клонуйте репозиторій**

   ```bash
   git clone https://github.com/Sigmanor/cherkasyoblenergo-api.git
   cd cherkasyoblenergo-api
   ```

3. **Налаштуйте змінні середовища**

   Створіть файл `.env` у кореневій директорії:

   ```properties
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=root
   DB_PASSWORD=your_strong_db_password
   DB_NAME=myCoolDB
   ADMIN_PASSWORD=your_strong_admin_password
   SERVER_PORT=3000
   ```

4. **Запустіть додаток**

   ```bash
   go run ./cmd/server/main.go
   ```

   Додаток автоматично створить необхідну базу даних при першому запуску.

### Збірка

Для збірки додатку для production:

```bash
go build -o cherkasyoblenergo_api ./cmd/server/main.go
```

## 🔑 Документація API

### Базовий URL

```
/cherkasyoblenergo/api
```

### Доступні ендпоінти

- `GET /blackout-schedule` - Отримати графіки відключень електроенергії
- `POST /api-keys` - Створити API ключ (адмін)
- `PATCH /api-keys` - Перегенерувати ключ або оновити ліміт (адмін)
- `DELETE /api-keys` - Видалити API ключ (адмін)

[Детальна документація API](API_UA.md)

## 🚦 Запуск тестів

Для запуску тестів локально:

```bash
go test ./...
```

## ⚡ Безкоштовний доступ до API

Оскільки цей API сервер постійно запущений у мене на хостингу (для моїх потреб) я можу надати вам обмежений доступ (2 запити/хв) до нього безкоштовно. Зв'яжіться зі мною за допомогою [email](mailto:dock-brunt-rarity@duck.com) для отримання доступу.
