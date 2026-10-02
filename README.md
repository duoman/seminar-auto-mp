# Marketplace — учебный проект для семинаров по КПО

Маркетплейс по продаже автомобилей от разных продавцов: дилеров, частников и аукционов.

---

## Зачем это нужно

**Проблема, которую решает проект.** В лекциях 1 и 2 мы разбирали DIP, DI и Composition Root. В лекции 3 — как писать полезные тесты. В лекции 4 — как применять паттерны GoF. Каждая тема по отдельности понятна, но у студентов часто нет **одного сквозного примера**, где всё это работает вместе.

**Решение.** Один проект — маркетплейс автомобилей — где:

- есть **интерфейсы и DI** (`ports/`, `cmd/api/main.go`),
- есть **unit-тесты** с разными видами двойников (`*_test.go`),
- есть **все четыре порождающих паттерна** (`listing/`),
- есть **реальный REST API и Postgres**, чтобы студенты видели production-подобную картину.

**Цель.** Студент после семинаров должен уметь:

- отличать unit-тест от интеграционного,
- писать тест по правилу Arrange-Act-Assert,
- выбирать подходящий тестовый двойник (fake, stub, spy, dummy),
- узнавать паттерн по проблеме, а не по названию,
- применять Factory Method, Builder, Prototype, Abstract Factory,
- видеть, где живут паттерны в реальном приложении.

---

## Что умеет

### Бизнес-логика

- Создание объявления о продаже автомобиля от **трёх типов продавцов**.
- Автоматическое формирование заголовка, описания, комиссии и статуса — **по-разному** для дилера, частника и аукциона.
- **Клонирование** объявления для дилеров, у которых много похожих машин.
- **Публикация** объявления с валидацией.
- Поддержка **трёх маркетплейсов**: standard, premium, auction — со своими комиссиями и правилами модерации.

### Инфраструктура

- **REST API** на стандартной библиотеке `net/http` (без фреймворков).
- **Postgres** через `database/sql` и `github.com/lib/pq`.
- **In-memory** хранилище для быстрого запуска и тестов.
- **Docker Compose** — Postgres + API поднимаются одной командой.
- **Миграции** — SQL-файл монтируется в `docker-entrypoint-initdb.d`.
- **Makefile** — все команды для сборки, тестов, запуска, Docker, БД.
- **Тесты** — unit и HTTP-интеграционные, с race detector.

---

## Паттерны, которые демонстрирует проект

| Паттерн | Файл | Проблема, которую решает |
|---|---|---|
| **Factory Method** | `listing/factory.go` | Выбор класса объявления в зависимости от типа продавца без `switch` в бизнес-коде |
| **Builder** | `listing/builder.go` | Сборка сложного объявления с десятком опциональных полей и валидацией |
| **Prototype** | `listing/prototype.go` | Клонирование похожих объявлений с глубокой копией фото и вложенных объектов |
| **Abstract Factory** | `listing/abstract_factory.go` | Создание семейства: объявление + комиссия + модерация — без смешивания |
| **Repository** (инфраструктурный) | `ports/repository.go`, `infra/*.go` | Отделить бизнес-логику от способа хранения |
| **Dependency Injection** | `cmd/api/main.go` | Передать зависимости явно, а не создавать их внутри |
| **Composition Root** | `cmd/api/main.go` | Единственное место, где выбираются конкретные реализации |

---

## Что освещается на семинаре A (юнит-тестирование)

### Темы

- Чем unit-тест отличается от интеграционного и e2e.
- Пирамида тестирования.
- Правило **Arrange-Act-Assert**.
- Стили тестов: **output-based**, **state-based**, **interaction-based**.
- Тестовые двойники: **dummy**, **fake**, **stub**, **spy**, **mock**.
- Хрупкий тест vs устойчивый к рефакторингу.
- Четыре критерия хорошего unit-теста.
- Что мокать, а что проверять настоящим объектом.
- Humble Object.
- DI и тестируемость.

### Файлы для разбора

| Файл | Что показывает |
|---|---|
| `listing/factory_test.go` | Arrange-Act-Assert, table-driven tests, подтесты через `t.Run` |
| `listing/builder_test.go` | Проверка ошибок через `errors.Is`, сентинел-ошибки |
| `listing/prototype_test.go` | Глубокая копия как способ изолировать тесты |
| `httpapi/handlers_test.go` | Интеграционные тесты HTTP без сети через `httptest` |

### Как запускать

```bash
make test          # все тесты
make test-v        # подробный вывод
make test-race     # с race detector
make test-cover    # покрытие → coverage.html
make test-service  # только listing
make test-http     # только httpapi
```

---

## Что освещается на семинаре B (паттерны GoF)

### Темы

- Паттерн — это имя для повторяющейся проблемы, а не готовый код.
- Порождающие, структурные, поведенческие паттерны — классификация GoF.
- **Factory Method** — выбор класса продукта отдельно от клиента.
- **Abstract Factory** — семейство связанных продуктов, которые нельзя смешивать.
- **Builder** — пошаговая сборка сложного объекта с валидацией.
- **Prototype** — клонирование, поверхностное и глубокое копирование.
- **Singleton** — почему в современной разработке его заменяет DI lifetime.
- Как выбирать паттерн по проблеме, а не по названию.
- Частые ошибки: YAGNI, «паттерн на будущее», смешивание ролей.

### Файлы для разбора

| Файл | Паттерн |
|---|---|
| `listing/factory.go` | Factory Method |
| `listing/builder.go` | Builder |
| `listing/prototype.go` | Prototype |
| `listing/abstract_factory.go` | Abstract Factory |
| `cmd/api/main.go` | DI + Composition Root |

### Как запускать

```bash
make cli           # CLI-демонстрация всех паттернов
make run           # REST API с in-memory хранилищем
make docker-up     # всё в Docker (Postgres + API)
```

---

## Быстрый старт

### 1. Клонировать и подтянуть зависимости

```bash
cd marketplace
make deps
```

Или вручную:

```bash
go mod tidy
```

### 2. Прогнать тесты

```bash
make test-v
```

### 3. Запустить CLI-демонстрацию паттернов

```bash
make cli
```

Ожидаемый вывод:

```
=== CLI-демонстрация паттернов ===

-- Factory Method --
[dealer  ] Tesla Model 3 2023, от официального дилера (комиссия 3%)
[private ] Tesla Model 3 2023, пробег 15000 км (комиссия 1%)
[auction ] Аукцион: Tesla Model 3 2023 (комиссия 5%)

-- Builder --
Tesla Model 3, идеальное состояние | фото: 2 | регион: СПб

-- Prototype --
original фото: 2
clone фото:    3
Оригинал не изменился: true

-- Abstract Factory --
[standard] комиссия 3%, модерация: false
[premium ] комиссия 5%, модерация: true
[auction ] комиссия 7%, модерация: true

Сохранено: 2
```

### 4. Запустить REST API локально

```bash
make run
```

Сервер на `localhost:8080`, хранилище — in-memory.

### 5. Запустить всё в Docker (Postgres + API)

```bash
make docker-up
```

API на `localhost:8080`, Postgres на `localhost:5432`.

Остановить:

```bash
make docker-down
```

### 6. Прогнать smoke-тест

В одном терминале:

```bash
make run
```

В другом:

```bash
make smoke
```

---

## Все команды Makefile

```bash
make help
```

Основные:

| Команда | Что делает |
|---|---|
| `make deps` | Скачать зависимости |
| `make build` | Собрать API и CLI в `./bin` |
| `make test` | Прогнать все тесты |
| `make test-race` | Тесты с race detector |
| `make test-cover` | Покрытие → `coverage.html` |
| `make run` | REST API локально (in-memory) |
| `make run-postgres` | REST API локально с Postgres |
| `make cli` | CLI-демонстрация паттернов |
| `make smoke` | curl-сценарий против `localhost:8080` |
| `make docker-up` | Всё в Docker в фоне |
| `make docker-down` | Остановить и удалить контейнеры |
| `make docker-logs` | Смотреть логи API |
| `make postgres-up` | Поднять только Postgres |
| `make postgres-wait` | Дождаться готовности Postgres |
| `make psql` | Открыть psql внутри контейнера |
| `make db-listings` | Показать объявления в БД |
| `make ci` | deps + fmt + vet + test-race |

Полный список — `make help`.

---

## REST API

База: `/api/v1`. Формат: JSON.

| Метод | Путь | Что делает | Паттерн |
|---|---|---|---|
| POST | `/api/v1/listings` | Создать объявление | Factory Method |
| POST | `/api/v1/listings/custom` | Кастомное объявление | Builder |
| POST | `/api/v1/listings/{id}/clone` | Клонировать | Prototype |
| POST | `/api/v1/listings/{id}/publish` | Опубликовать | — |
| GET | `/api/v1/listings` | Список объявлений | — |
| GET | `/api/v1/listings/{id}` | Одно объявление | — |
| GET | `/api/v1/marketplaces` | Тарифы маркетплейсов | Abstract Factory |
| GET | `/healthz` | Health check | — |

### Примеры запросов

**Создать объявление от дилера:**

```bash
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
```

**Список:**

```bash
curl -s localhost:8080/api/v1/listings | jq
```

**Клонировать:**

```bash
curl -s -X POST localhost:8080/api/v1/listings/L-D-1-CAR-1/clone | jq
```

**Опубликовать:**

```bash
curl -s -X POST localhost:8080/api/v1/listings/L-D-1-CAR-1/publish | jq
```

**Тарифы:**

```bash
curl -s localhost:8080/api/v1/marketplaces | jq
```

---

## Структура проекта

```
marketplace/
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── .env.example
├── .dockerignore
├── README.md
│
├── domain/                 # предметная область
│   ├── car.go              # Car
│   ├── seller.go           # Seller, SellerType
│   └── listing.go          # Listing, ListingStatus
│
├── ports/                  # интерфейсы (DIP)
│   └── repository.go       # ListingRepository
│
├── listing/                # паттерны + их тесты
│   ├── factory.go          # Factory Method
│   ├── factory_test.go
│   ├── builder.go          # Builder
│   ├── builder_test.go
│   ├── prototype.go        # Prototype
│   ├── prototype_test.go
│   ├── abstract_factory.go # Abstract Factory
│   └── abstract_factory_test.go
│
├── infra/                  # реализации хранилищ
│   ├── repo_memory.go      # in-memory
│   └── repo_postgres.go    # Postgres
│
├── httpapi/                # REST-слой
│   ├── handlers.go
│   └── handlers_test.go
│
├── migrations/             # миграции БД
│   └── 001_init.sql
│
└── cmd/                    # точки входа
    ├── api/                # REST API
    │   └── main.go
    └── marketplace/        # CLI-демонстрация
        └── main.go
```

### Слои

| Слой | Пакет | Отвечает за |
|---|---|---|
| Домен | `domain/` | Чистые типы предметной области |
| Порты | `ports/` | Интерфейсы, которые описывает потребитель |
| Приложение | `listing/` | Бизнес-логика и паттерны |
| Инфраструктура | `infra/` | Реализации хранилищ, адаптеры |
| Транспорт | `httpapi/` | HTTP-хендлеры, JSON |
| Точка входа | `cmd/` | Composition Root |

**Направление зависимостей:**

```
cmd → httpapi → listing → ports ← infra
                  ↓
                domain
```

Бизнес-логика (`listing/`) зависит только от `domain` и `ports`. Инфраструктура (`infra/`) — тоже. `httpapi/` знает про бизнес-логику, но не знает про способ хранения. Выбор реализации — только в `cmd/`.

---

## Сценарии на семинаре

### Семинар A — тесты

**Открывать по порядку:**

1. `listing/factory_test.go` — показать Arrange-Act-Assert и table-driven.
2. `listing/builder_test.go` — показать проверку ошибок через `errors.Is`.
3. `listing/prototype_test.go` — показать, как глубокая копия изолирует тесты.
4. `httpapi/handlers_test.go` — показать `httptest`, тесты без сети и БД.

**Запускать:**

```bash
make test-v
make test-race
make test-cover
```

Открыть `coverage.html` в браузере — показать покрытие по слоям.

### Семинар B — паттерны

**Открывать по порядку:**

1. `listing/factory.go` — Factory Method: почему не `switch` в бизнес-коде.
2. `listing/builder.go` — Builder: почему не телескопический конструктор.
3. `listing/prototype.go` — Prototype: почему не ручное копирование.
4. `listing/abstract_factory.go` — Abstract Factory: почему семейства нельзя смешивать.
5. `cmd/api/main.go` — DI + Composition Root: где выбирается хранилище.

**Запускать:**

```bash
make cli                # CLI-демонстрация
make docker-up          # всё в Docker
make db-listings        # проверить БД
make docker-down
```

---

## Технологии

- **Go 1.22** — язык.
- **`net/http`** — HTTP-сервер из стандартной библиотеки.
- **`database/sql`** + **`github.com/lib/pq`** — работа с Postgres.
- **`testing`** + **`net/http/httptest`** — тесты.
- **Postgres 16** — БД.
- **Docker + Docker Compose** — инфраструктура.
- **Make** — управление командами.

**Внешних зависимостей только одна** — драйвер Postgres. Это принципиально: студенты видят, что бизнес-логика, паттерны и тесты не требуют фреймворков.

---

## Переменные окружения

| Переменная | Значения | По умолчанию |
|---|---|---|
| `STORAGE` | `memory` \| `postgres` | `memory` |
| `DATABASE_URL` | строка подключения | — |
| `PORT` | номер порта | `8080` |

Пример:

```bash
STORAGE=postgres \
DATABASE_URL='postgres://marketplace:marketplace@localhost:5432/marketplace?sslmode=disable' \
PORT=9090 \
make run
```

---

## Как переключиться между хранилищами

**In-memory** — по умолчанию:

```bash
make run
```

**Postgres** — через Docker:

```bash
make postgres-up
make postgres-wait
make run-postgres
```

**Postgres + API целиком в Docker:**

```bash
make docker-up
```

Один и тот же код, одно и то же API. Меняется только переменная `STORAGE`. Это иллюстрация DIP в действии: бизнес-логика не знает, куда пишутся данные.

---

## Отладка

**API не стартует с Postgres** — дождись готовности БД:

```bash
make postgres-up
make postgres-wait
```

**Миграции не применились** — Postgres выполняет их только при первом старте. Если схема менялась:

```bash
make docker-clean     # сбросит БД и volume
make docker-up
```

**Порт 8080 занят** — переопредели:

```bash
make run PORT=9090
make smoke PORT=9090
```

**Тесты с race detector долгие** — используй только для семинара A:

```bash
make test          # без -race
make test-race     # с -race, дольше
```

---

## Что дальше

Проект покрывает лекции 3 (юнит-тестирование) и 4 (паттерны GoF). В следующих лекциях курса появятся:

- **Структурные паттерны** — Adapter, Decorator, Facade, Composite.
- **Поведенческие паттерны** — Strategy, State, Observer, Template Method.
- **CI/CD** — тесты как часть pipeline от commit до production.

Проект легко расширяется под эти темы: `listing/` — для поведенческих, `infra/` — для структурных, `httpapi/` — для интеграционных сценариев.

---

## Лицензия

Учебный проект. Используется на семинарах курса «Конструирование программного обеспечения».