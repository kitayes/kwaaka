# Kwaaka

Pet-проект по ТЗ: система для парсинга меню из Google Sheets и управления статусами товаров.  
Архитектура разделена на два сервиса:

- **API** — принимает HTTP-запросы от клиента.
- **Worker** — асинхронный воркер, который:
    - парсит меню из Google Sheets,
    - сохраняет его в MongoDB,
    - обрабатывает события изменения статуса товара.

Оба сервиса общаются через **RabbitMQ** и используют общую **MongoDB**.

---

## Архитектура

### Компоненты

- **API сервис (`api/`)**
    - HTTP API (`/api/v1/...`)
    - Создание задач парсинга меню.
    - Получение статуса задачи.
    - Получение меню.
    - Изменение статуса товара (через публикацию события в очередь).

- **Worker сервис (`worker/`)**
    - Подписчик на очередь задач парсинга меню.
    - Подписчик на очередь событий статуса товара.
    - Парсинг Google Sheets → создание меню в MongoDB.
    - Обновление статуса товара в меню + запись аудита в `product_status_audit`.
    - Health-эндпоинт `/healthz`.

- **MongoDB**
    - `menus` — меню ресторана.
    - `parsing_tasks` — задачи на парсинг.
    - `product_status_audit` — аудит изменений статуса товаров.

- **RabbitMQ**
    - Очередь для парсинга меню (например, `menu-parsing`).
    - Очередь для изменений статуса товаров (например, `product-status`).

---

## Стек

- Go 1.24
- Gin
- MongoDB 7 (mongo-driver)
- RabbitMQ 3.13 (amqp)
- Google Sheets API (`google.golang.org/api/sheets/v4`)
- Docker + docker-compose

---

## Рекомендуемый способ запуска

Рекомендуемый способ — запуск всего стека через **Docker + docker-compose**.

### 1. Предварительные требования

- Установлен **Docker** и **docker-compose**
- Есть JSON-ключ сервисного аккаунта Google с доступом к нужным Google Sheets

Структура проекта (важные части):

```text
kwaaka/
  api/
    .env
    Dockerfile
  worker/
    .env
    Dockerfile
  infrastructure/
    docker-compose.yml
  credentials/
    google-service-account.json   # ключ сервисного аккаунта (не коммитить в git)
