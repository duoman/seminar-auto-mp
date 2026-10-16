# Marketplace + Delivery

Монорепозиторий учебного проекта для курса КПО. Два независимых сервиса,
которые общаются по HTTP и вместе демонстрируют материал лекций 1–6.

---

## 📖 Оглавление

- [Что это](#-что-это)
- [Быстрый старт](#-быстрый-старт)
- [Архитектура](#-архитектура)
- [Бизнес-логика маркетплейса](#-бизнес-логика-маркетплейса)
- [Бизнес-логика доставки](#-бизнес-логика-доставки)
- [Как сервисы взаимодействуют](#-как-сервисы-взаимодействуют)
- [Что демонстрирует проект](#-что-демонстрирует-проект)
- [API](#-api)
- [Структура проекта](#-структура-проекта)
- [Команды Makefile](#-команды-makefile)
- [Переменные окружения](#-переменные-окружения)
- [Проверка](#-проверка)
- [Для семинаров](#-для-семинаров)

---

## 🎯 Что это

**Marketplace** — маркетплейс по продаже автомобилей. Продавцы публикуют объявления,
покупатели покупают машины, платформа берёт комиссию.

**Delivery** — микросервис доставки. Когда машина продана, marketplace вызывает
delivery, чтобы запланировать перевозку от продавца к покупателю.

**Зачем два сервиса.** Один сервис не может показать все паттерны из лекций 4–6
органично. Marketplace — домен машин и объявлений, delivery — домен маршрутов и
доставки. Вместе они покрывают весь материал курса и дают естественный повод
для микросервисного взаимодействия.

**Ключевая идея.** Каждый паттерн появляется **там, где была конкретная боль** —
не «на будущее». Если бы боли не было, паттерн был бы лишним. Это YAGNI.

---

## 🚀 Быстрый старт

### Всё в Docker

```bash
make up        # поднять postgres + delivery + marketplace
make smoke     # прогнать сценарий продажи + доставки
make down      # остановить
make clean     # остановить + сбросить БД
```

### Локально (два терминала)

```bash
# терминал 1 — delivery на :8081
make delivery

# терминал 2 — marketplace на :8080
make marketplace
```

### Каждый сервис отдельно

```bash
cd marketplace && go run ./cmd/api
cd delivery    && go run ./cmd/api
```

Каждый сервис самодостаточен: свой `go.mod`, свой `Dockerfile`.

---

## 🏛 Архитектура

```
┌─────────────────────────┐         HTTP          ┌─────────────────────────┐
│      marketplace        │  ──────────────────►  │       delivery          │
│        :8080            │                       │        :8081            │
│                         │  ◄──────────────────  │                         │
│  Domain:                │    delivery_id        │  Domain:                │
│   - Car                 │                       │   - Delivery            │
│   - Seller              │                       │   - Route               │
│   - Listing             │                       │   - Point               │
│                         │                       │                         │
│  Patterns:              │                       │  Patterns:              │
│   - Factory Method      │                       │   - Strategy            │
│   - Builder             │                       │   - State               │
│   - Prototype           │                       │   - Template Method     │
│   - Abstract Factory    │                       │   - Observer            │
│   - Repository          │                       │   - Command             │
│   - DIP / DI            │                       │   - Visitor             │
│                         │                       │   - Adapter             │
│                         │                       │   - Proxy               │
│  Storage:               │                       │   - Decorator           │
│   - in-memory           │                       │   - Flyweight           │
│   - PostgreSQL          │                       │   - Facade              │
└─────────────────────────┘                       └─────────────────────────┘
```

**Направление зависимостей внутри каждого сервиса:**

```
cmd → httpapi → listing/behavior → ports ← infra
                     ↓
                  domain
```

Бизнес-логика не знает про инфраструктуру. Реализации зависят от интерфейсов.
Composition Root (`cmd/api/main.go`) — единственное место выбора.

---

## 🚗 Бизнес-логика маркетплейса

### Роли

| Роль | Что делает |
|---|---|
| **Дилер** | Публикует машины с гарантией и тест-драйвом |
| **Частник** | Продаёт личную машину |
| **Аукцион** | Выставляет лот на торги |
| **Покупатель** | Покупает машину, получает доставку |
| **Администратор** | Управляет тарифами |

### Сценарий 1. Создание объявления

1. Продавец выбирает машину: `Tesla Model 3, 2023, 15 000 км, 5 000 000 ₽`.
2. Указывает тип продавца: `dealer`, `private` или `auction`.
3. Система автоматически формирует объявление по правилам типа:
   - **Dealer** → заголовок «от официального дилера», гарантия ✅, тест-драйв ✅, комиссия 3%
   - **Private** → заголовок с пробегом, гарантии нет, тест-драйв ✅, комиссия 1%
   - **Auction** → «Аукцион: …», статус `draft`, комиссия 5%
4. Объявление сохраняется в статусе `pending` (на модерацию).

### Сценарий 2. Публикация

1. Модератор проверяет объявление.
2. `POST /listings/{id}/publish` → статус `active`.
3. Покупатели видят объявление в каталоге.

### Сценарий 3. Кастомное объявление (Builder)

Дилер хочет нестандартное оформление:

1. Заполняет заголовок, описание, регион, фото, комиссию.
2. Флаги `has_warranty`, `has_test_drive`.
3. `Build()` валидирует обязательные поля.
4. Готовое объявление получает `pending`.

### Сценарий 4. Клонирование (Prototype)

У дилера 10 одинаковых машин:

1. Создаёт одно объявление.
2. `POST /listings/{id}/clone` × 9.
3. Меняет только VIN и цену.
4. Важно: фото копируются глубоко — изменение клона не влияет на оригинал.

### Сценарий 5. Продажа

1. Покупатель выбрал машину → `POST /listings/{id}/sell`.
2. Marketplace проверяет, что объявление `active`.
3. Вызывает delivery → `POST /api/v1/deliveries`.
4. Получает `delivery_id`.
5. Меняет статус объявления на `sold`, сохраняет `delivery_id`.

### Правила модерации по тарифам

| Тариф | Комиссия | Модерация |
|---|---|---|
| **Standard** | 3% | автоматическая |
| **Premium** | 5% | ручная |
| **Auction** | 7% | ручная |

Тариф определяет семейство правил: комиссию, модерацию, тип объявления.
Смешивать нельзя — объявление от дилера не может иметь правил аукциона.

---

## 🚚 Бизнес-логика доставки

### Типы доставки

| Тип | Алгоритм | Цена | Пайплайн |
|---|---|---|---|
| **Express** | Самый быстрый | 15 000 ₽ | City (SMS) |
| **Economy** | Самый дешёвый | 8 000 ₽ | City (SMS) |
| **Eco** | Эко-маршрут | 12 000 ₽ | Intercity (push) |
| **International** | Международный | 60 000 ₽ | Intl (email) |

### Жизненный цикл доставки

```
created ──schedule──► scheduled ──start──► in_transit ──complete──► delivered
   │                       │                    │
   └───────cancel──────────┴────────────────────┘
                                                  ✗ нельзя из delivered
```

Недопустимые переходы возвращают ошибку — например, нельзя начать доставку,
которая ещё не запланирована.

### Сценарий 1. Создание

`POST /api/v1/deliveries` с телом:

```json
{
  "listing_id": "L-D-1-CAR-1",
  "customer_id": "C-1",
  "type": "express",
  "from_city": "Moscow",
  "to_city": "SPb"
}
```

Что происходит внутри:

1. **Facade** координирует сценарий.
2. **Strategy** выбирает алгоритм по типу доставки.
3. **Flyweight** достаёт шаблон маршрута `Moscow|SPb` (создаётся один раз).
4. **Proxy** → **Adapter** обращается к Maps SDK (лениво).
5. **State** переводит `created → scheduled`.
6. **Command** оборачивает смену статуса (можно откатить).
7. **Observer** уведомляет подписчиков: customer, warehouse, analytics, marketplace.

### Сценарий 2. Запуск

`POST /deliveries/{id}/start`:

1. **State**: `scheduled → in_transit`.
2. **Template Method**: прогон пайплайна `Pickup → Transport → Dropoff → Notify`.
   Порядок шагов фиксирован, реализация зависит от типа.
3. **Observer**: уведомления.

### Сценарий 3. Завершение

`POST /deliveries/{id}/complete`:

1. **State**: `in_transit → delivered`.
2. **Observer**: уведомления.

### Сценарий 4. Undo

`POST /deliveries/{id}/undo`:

**Command** откатывает последнюю смену статуса.

### Сценарий 5. Отчёты

`GET /api/v1/reports`:

**Visitor** обходит все доставки двумя операциями:

- **InsuranceVisitor** — считает страховку (разная для local / long-haul / international).
- **CostReportVisitor** — собирает строки отчёта.

Ответ:

```json
{
  "total_insurance": 2610,
  "report": [
    "Local D-1: 1000 RUB",
    "Long-haul D-2: 5000 RUB",
    "International D-3: 50000 RUB"
  ]
}
```

---

## 🔗 Как сервисы взаимодействуют

**Marketplace → Delivery:**

```
POST /listings/{id}/sell
       │
       ▼
marketplace.httpapi.Server.SellListing
       │
       ▼
ports.DeliveryClient.Schedule(listingID, customerID, from, to, type)
       │
       ▼
infra.HTTPDeliveryClient — POST http://delivery:8081/api/v1/deliveries
       │
       ▼
delivery.httpapi.Server.Create
       │
       ▼
structural.ScheduleFacade.Schedule
       │
       ▼
Strategy + Flyweight + Proxy + State + Command + Observer
```

**Что важно для семинара:**

- Marketplace не знает, как устроена delivery. Он знает только интерфейс
  `DeliveryClient` — это **DIP**.
- Конкретная реализация (`HTTPDeliveryClient`) собирается в Composition Root —
  это **DI**.
- Если delivery переедет в Kubernetes, поменяется только URL. Код не тронем.

**Ответ delivery в marketplace:**

```json
{ "id": "DEL-L-D-1-CAR-1", "status": "scheduled" }
```

`delivery_id` сохраняется в объявлении — можно потом посмотреть статус доставки.

---

## 📚 Что демонстрирует проект

Каждый паттерн появился там, где была конкретная боль. Не «на будущее».

### Лекция 1 — Введение в КПО

**Где:** `marketplace/domain/`, `marketplace/listing/`

| Идея | Где |
|---|---|
| Интерфейсы как контракт | `domain.Engine`, `ports.ListingRepository` |
| Композиция | `Listing.Car` владеет `*Car` |
| Агрегация | `Listing.Seller` — внешний объект |
| SOLID / SRP | Разделение `Car`, `Listing`, `Seller` |
| SOLID / OCP | `listing.FactoryFor` — новый тип продавца без правки клиента |
| SOLID / DIP | Бизнес-логика зависит от `ports`, не от `infra` |

### Лекция 2 — Внедрение зависимостей

**Где:** `marketplace/cmd/api/main.go`, `delivery/cmd/api/main.go`

| Идея | Где |
|---|---|
| DIP | `httpapi.Server` зависит от `ports.ListingRepository` |
| Constructor injection | `NewServer(repo, delivery)` |
| Composition Root | `pickRepository()`, `pickDeliveryClient()` |
| Выбор storage | `STORAGE=memory` vs `STORAGE=postgres` |
| Выбор delivery URL | `DELIVERY_URL` — runtime |

### Лекция 3 — Юнит-тестирование

**Где:** `*_test.go` в обоих сервисах (задел на семинар)

- `fakeCarRepo`, `stubPayment`, `spyNotifier` — тестовые двойники.
- Arrange-Act-Assert.
- Table-driven tests.
- Глубокая копия в `Prototype` как способ изоляции тестов.

### Лекция 4 — Порождающие паттерны GoF

**Где:** `marketplace/listing/`

| Паттерн | Файл | Проблема |
|---|---|---|
| **Factory Method** | `factory.go` | Выбор класса объявления по типу продавца |
| **Builder** | `builder.go` | Сложное объявление с опциональными полями |
| **Prototype** | `prototype.go` | Клонирование похожих объявлений с глубокой копией |
| **Abstract Factory** | `abstract_factory.go` | Семейства Standard / Premium / Auction |

### Лекция 5 — Поведенческие паттерны GoF

**Где:** `delivery/behavior/`

| Паттерн | Файл | Проблема |
|---|---|---|
| **Strategy** | `strategy.go` | Алгоритм расчёта маршрута по типу доставки |
| **State** | `state.go` | Жизненный цикл доставки `created → delivered` |
| **Template Method** | `template_method.go` | Каркас пайплайна, разные шаги по типу |
| **Observer** | `observer.go` | Уведомления подписчиков при смене статуса |
| **Command** | `command.go` | Смена статуса с undo |
| **Visitor** | `visitor.go` | Страховка и отчёт над иерархией доставок |

### Лекция 6 — Структурные паттерны GoF

**Где:** `delivery/structural/`

| Паттерн | Файл | Проблема |
|---|---|---|
| **Adapter** | `adapter.go` | Сторонний Maps SDK → наш `GeoProvider` |
| **Proxy** | `proxy.go` | Ленивое создание SDK |
| **Decorator** | `decorator.go` | Логирование + кеш вокруг калькулятора |
| **Flyweight** | `flyweight.go` | Шаблоны маршрутов — не дублируем |
| **Facade** | `facade.go` | Один вход в планирование доставки |

---

## 🌐 API

### Marketplace (`:8080`)

| Метод | Путь | Что делает | Паттерн |
|---|---|---|---|
| POST | `/api/v1/listings` | Создать объявление | Factory Method |
| POST | `/api/v1/listings/custom` | Кастомное объявление | Builder |
| POST | `/api/v1/listings/{id}/clone` | Клонировать | Prototype |
| POST | `/api/v1/listings/{id}/publish` | Опубликовать | — |
| POST | `/api/v1/listings/{id}/sell` | Продать (вызывает delivery) | DIP + DI |
| GET | `/api/v1/listings` | Список | — |
| GET | `/api/v1/listings/{id}` | Одно | — |
| GET | `/api/v1/marketplaces` | Тарифы | Abstract Factory |
| GET | `/healthz` | Health | — |

### Delivery (`:8081`)

| Метод | Путь | Что делает | Паттерн |
|---|---|---|---|
| POST | `/api/v1/deliveries` | Создать доставку | Facade + Strategy + Flyweight |
| POST | `/api/v1/deliveries/{id}/start` | Запустить | State + Template Method |
| POST | `/api/v1/deliveries/{id}/complete` | Завершить | State |
| POST | `/api/v1/deliveries/{id}/undo` | Откатить | Command |
| GET | `/api/v1/deliveries` | Список | — |
| GET | `/api/v1/deliveries/{id}` | Одна | — |
| GET | `/api/v1/reports` | Отчёты | Visitor |
| GET | `/healthz` | Health | — |

---

## 📁 Структура проекта

```
.
├── README.md
├── Makefile
├── docker-compose.yml
│
├── marketplace/                    # :8080
│   ├── go.mod
│   ├── Dockerfile
│   ├── domain/                     # Car, Seller, Listing
│   │   ├── car.go
│   │   ├── seller.go
│   │   └── listing.go
│   ├── listing/                    # Порождающие паттерны
│   │   ├── factory.go              # Factory Method
│   │   ├── builder.go              # Builder
│   │   ├── prototype.go            # Prototype
│   │   └── abstract_factory.go     # Abstract Factory
│   ├── ports/                      # Интерфейсы (DIP)
│   │   ├── repository.go
│   │   └── delivery.go
│   ├── infra/                      # Реализации
│   │   ├── repo_memory.go
│   │   ├── repo_postgres.go
│   │   └── delivery_http.go        # HTTP-клиент к delivery
│   ├── httpapi/
│   │   └── handlers.go
│   ├── migrations/
│   │   └── 001_init.sql
│   └── cmd/api/
│       └── main.go                 # Composition Root
│
└── delivery/                       # :8081
    ├── go.mod
    ├── Dockerfile
    ├── domain/
    │   └── delivery.go             # Delivery, Route, Point
    ├── behavior/                   # Поведенческие паттерны
    │   ├── strategy.go             # Strategy
    │   ├── state.go                # State
    │   ├── template_method.go      # Template Method
    │   ├── observer.go             # Observer
    │   ├── command.go              # Command
    │   └── visitor.go              # Visitor
    ├── structural/                 # Структурные паттерны
    │   ├── adapter.go              # Adapter
    │   ├── proxy.go                # Proxy
    │   ├── decorator.go            # Decorator
    │   ├── flyweight.go            # Flyweight
    │   └── facade.go               # Facade
    ├── ports/
    │   └── repository.go
    ├── infra/
    │   └── repo_memory.go
    ├── httpapi/
    │   └── handlers.go
    └── cmd/api/
        └── main.go                 # Composition Root
```

---

## 🛠 Команды Makefile

```bash
make help          # список всех команд
```

### Docker

| Команда | Что делает |
|---|---|
| `make up` | Поднять всё (postgres + delivery + marketplace) |
| `make down` | Остановить |
| `make clean` | Остановить + удалить volume БД |
| `make logs` | Логи |
| `make ps` | Контейнеры |

### Локально

| Команда | Что делает |
|---|---|
| `make delivery` | Delivery локально (:8081) |
| `make marketplace` | Marketplace локально, in-memory (:8080) |
| `make marketplace-postgres` | Marketplace локально + Postgres |

### Разработка

| Команда | Что делает |
|---|---|
| `make build` | Собрать оба сервиса |
| `make test` | Тесты обоих |
| `make tidy` | go mod tidy |
| `make fmt` | go fmt |

### Postgres

| Команда | Что делает |
|---|---|
| `make postgres-up` | Поднять только Postgres |
| `make psql` | Открыть psql |
| `make db-listings` | Показать объявления |

### Smoke

| Команда | Что делает |
|---|---|
| `make smoke` | Полный сценарий: создать → опубликовать → продать → доставить |

---

## ⚙️ Переменные окружения

### Marketplace

| Переменная | Значения | По умолчанию |
|---|---|---|
| `PORT` | Номер порта | `8080` |
| `STORAGE` | `memory` \| `postgres` | `memory` |
| `DATABASE_URL` | Строка подключения | — |
| `DELIVERY_URL` | URL delivery | `http://localhost:8081` |

### Delivery

| Переменная | Значения | По умолчанию |
|---|---|---|
| `PORT` | Номер порта | `8081` |

---

## ✅ Проверка

### Health

```bash
curl -s localhost:8080/healthz | jq
curl -s localhost:8081/healthz | jq
```

### Полный сценарий

```bash
# 1. Создать объявление от дилера
curl -s -X POST localhost:8080/api/v1/listings \
  -H "Content-Type: application/json" \
  -d '{
    "car_id": "CAR-1",
    "brand": "Tesla",
    "model": "Model 3",
    "year": 2023,
    "mileage": 15000,
    "price": 5000000,
    "seller_id": "D-1",
    "seller_name": "AutoDealer",
    "seller_type": "dealer"
  }' | jq

# 2. Опубликовать
curl -s -X POST localhost:8080/api/v1/listings/L-D-1-CAR-1/publish | jq

# 3. Продать — marketplace вызовет delivery
curl -s -X POST localhost:8080/api/v1/listings/L-D-1-CAR-1/sell | jq

# 4. Проверить, что доставка создана
curl -s localhost:8081/api/v1/deliveries | jq

# 5. Запустить доставку
curl -s -X POST localhost:8081/api/v1/deliveries/DEL-L-D-1-CAR-1/start | jq

# 6. Завершить
curl -s -X POST localhost:8081/api/v1/deliveries/DEL-L-D-1-CAR-1/complete | jq

# 7. Отчёты (Visitor)
curl -s localhost:8081/api/v1/reports | jq
```

### Или одной командой

```bash
make smoke
```

---

## 📝 Для семинаров

**Семинар A — юнит-тестирование.** Открывать `*_test.go`, запускать `make test`.

**Семинар B — паттерны GoF.** Открывать `listing/`, `behavior/`, `structural/`.
Запускать `make marketplace` и `make delivery`.

**Семинар по микросервисам.** Показать `docker-compose.yml`, `HTTPDeliveryClient`,
объяснить DIP между сервисами.

---
