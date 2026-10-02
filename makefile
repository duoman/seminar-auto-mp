# ============================================================
# Marketplace — Makefile
# ============================================================
#
# Основные команды:
#   make help          — список всех команд
#   make test          — прогнать все тесты
#   make run           — REST API локально (in-memory)
#   make docker-up     — всё в Docker (Postgres + API)
#   make docker-down   — остановить Docker
#   make smoke         — прогнать curl-сценарий
#
# ============================================================

# ---------- настройки ----------
APP_NAME        := marketplace
API_BIN         := bin/api
CLI_BIN         := bin/cli

GO              ?= go
DOCKER          ?= docker
DOCKER_COMPOSE  ?= $(DOCKER) compose

# Параметры Postgres для локального запуска
PG_USER         ?= marketplace
PG_PASSWORD     ?= marketplace
PG_DB           ?= marketplace
PG_HOST         ?= localhost
PG_PORT         ?= 5432
DATABASE_URL    ?= postgres://$(PG_USER):$(PG_PASSWORD)@$(PG_HOST):$(PG_PORT)/$(PG_DB)?sslmode=disable

PORT            ?= 8080

# Цветной вывод
GREEN  := \033[0;32m
YELLOW := \033[0;33m
RED    := \033[0;31m
RESET  := \033[0m

.PHONY: help
help: ## Показать список команд
	@echo "$(GREEN)Marketplace — доступные команды:$(RESET)"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-18s$(RESET) %s\n", $$1, $$2}'

# ============================================================
# Сборка и запуск локально
# ============================================================

.PHONY: deps
deps: ## Скачать зависимости
	$(GO) mod tidy
	$(GO) mod download

.PHONY: build
build: deps ## Собрать API и CLI в ./bin
	@mkdir -p bin
	$(GO) build -o $(API_BIN) ./cmd/api
	$(GO) build -o $(CLI_BIN) ./cmd/marketplace
	@echo "$(GREEN)✓ собрано: $(API_BIN), $(CLI_BIN)$(RESET)"

.PHONY: build-api
build-api: deps ## Собрать только API
	@mkdir -p bin
	$(GO) build -o $(API_BIN) ./cmd/api
	@echo "$(GREEN)✓ $(API_BIN)$(RESET)"

.PHONY: build-cli
build-cli: deps ## Собрать только CLI
	@mkdir -p bin
	$(GO) build -o $(CLI_BIN) ./cmd/marketplace
	@echo "$(GREEN)✓ $(CLI_BIN)$(RESET)"

.PHONY: run
run: ## Запустить API локально (in-memory, STORAGE=memory)
	STORAGE=memory PORT=$(PORT) $(GO) run ./cmd/api

.PHONY: run-postgres
run-postgres: ## Запустить API локально с Postgres (нужен `make postgres-up`)
	STORAGE=postgres DATABASE_URL='$(DATABASE_URL)' PORT=$(PORT) $(GO) run ./cmd/api

.PHONY: cli
cli: ## Запустить CLI-демонстрацию паттернов
	$(GO) run ./cmd/marketplace

# ============================================================
# Тесты
# ============================================================

.PHONY: test
test: ## Прогнать все тесты
	$(GO) test ./...

.PHONY: test-v
test-v: ## Прогнать тесты с подробным выводом
	$(GO) test -v ./...

.PHONY: test-race
test-race: ## Прогнать тесты с race detector
	$(GO) test -race -v ./...

.PHONY: test-cover
test-cover: ## Тесты с покрытием (создаст coverage.html)
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ покрытие: coverage.html$(RESET)"

.PHONY: test-service
test-service: ## Тесты только listing-паттернов
	$(GO) test -v ./listing/...

.PHONY: test-http
test-http: ## Тесты только HTTP-слоя
	$(GO) test -v ./httpapi/...

# ============================================================
# Линт и форматирование
# ============================================================

.PHONY: fmt
fmt: ## Форматировать код
	$(GO) fmt ./...
	@echo "$(GREEN)✓ форматирование готово$(RESET)"

.PHONY: vet
vet: ## Прогнать go vet
	$(GO) vet ./...
	@echo "$(GREEN)✓ vet пройден$(RESET)"

.PHONY: tidy
tidy: ## Привести go.mod в порядок
	$(GO) mod tidy

# ============================================================
# Docker: всё в контейнерах
# ============================================================

.PHONY: docker-build
docker-build: ## Собрать Docker-образы
	$(DOCKER_COMPOSE) build
	@echo "$(GREEN)✓ образы собраны$(RESET)"

.PHONY: docker-up
docker-up: ## Поднять всё в Docker (Postgres + API) в фоне
	$(DOCKER_COMPOSE) up --build -d
	@echo "$(GREEN)✓ запущено. API: http://localhost:$(PORT)$(RESET)"
	@echo "$(YELLOW)  логи: make docker-logs$(RESET)"

.PHONY: docker-down
docker-down: ## Остановить и удалить контейнеры
	$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✓ остановлено$(RESET)"

.PHONY: docker-stop
docker-stop: ## Остановить контейнеры (без удаления)
	$(DOCKER_COMPOSE) stop
	@echo "$(GREEN)✓ остановлено$(RESET)"

.PHONY: docker-restart
docker-restart: docker-down docker-up ## Перезапустить всё в Docker

.PHONY: docker-logs
docker-logs: ## Смотреть логи API
	$(DOCKER_COMPOSE) logs -f api

.PHONY: docker-logs-db
docker-logs-db: ## Смотреть логи Postgres
	$(DOCKER_COMPOSE) logs -f postgres

.PHONY: docker-ps
docker-ps: ## Показать запущенные контейнеры
	$(DOCKER_COMPOSE) ps

.PHONY: docker-clean
docker-clean: ## Удалить контейнеры + тома (сбросит БД!)
	$(DOCKER_COMPOSE) down -v
	@echo "$(YELLOW)⚠ удалены контейнеры и volume Postgres$(RESET)"

# ============================================================
# Только Postgres (для локального запуска)
# ============================================================

.PHONY: postgres-up
postgres-up: ## Поднять только Postgres в Docker
	$(DOCKER_COMPOSE) up -d postgres
	@echo "$(GREEN)✓ Postgres на localhost:$(PG_PORT)$(RESET)"
	@echo "$(YELLOW)  Дождись ready: make postgres-wait$(RESET)"

.PHONY: postgres-down
postgres-down: ## Остановить только Postgres
	$(DOCKER_COMPOSE) stop postgres

.PHONY: postgres-wait
postgres-wait: ## Дождаться готовности Postgres
	@echo "Ждём Postgres..."
	@until $(DOCKER) exec marketplace-postgres pg_isready -U $(PG_USER) -d $(PG_DB) >/dev/null 2>&1; do \
		sleep 1; \
	done
	@echo "$(GREEN)✓ Postgres готов$(RESET)"

.PHONY: psql
psql: ## Открыть psql внутри контейнера
	$(DOCKER) exec -it marketplace-postgres psql -U $(PG_USER) -d $(PG_DB)

.PHONY: db-listings
db-listings: ## Показать все объявления в БД
	$(DOCKER) exec -it marketplace-postgres psql -U $(PG_USER) -d $(PG_DB) \
		-c "SELECT id, title, status, commission FROM listings ORDER BY created_at DESC;"

.PHONY: db-truncate
db-truncate: ## Очистить таблицу listings
	$(DOCKER) exec -it marketplace-postgres psql -U $(PG_USER) -d $(PG_DB) \
		-c "TRUNCATE listings;"
	@echo "$(GREEN)✓ таблица очищена$(RESET)"

# ============================================================
# Smoke-тест: curl-сценарий
# ============================================================

.PHONY: smoke
smoke: ## Прогнать curl-сценарий против localhost:$(PORT)
	@echo "$(GREEN)=== Smoke test ===$(RESET)"
	@echo "1. Health check"
	@curl -sf localhost:$(PORT)/healthz | jq . || (echo "$(RED)✗ API недоступен$(RESET)"; exit 1)
	@echo ""
	@echo "2. Создать объявление от дилера"
	@curl -sf -X POST localhost:$(PORT)/api/v1/listings \
		-H "Content-Type: application/json" \
		-d '{"car_id":"CAR-1","brand":"Tesla","model":"Model 3","year":2023,"mileage":15000,"price":5000000,"seller_id":"D-1","seller_name":"AutoDealer","seller_type":"dealer"}' | jq .
	@echo ""
	@echo "3. Список объявлений"
	@curl -sf localhost:$(PORT)/api/v1/listings | jq .
	@echo ""
	@echo "4. Клонировать"
	@curl -sf -X POST localhost:$(PORT)/api/v1/listings/L-D-1-CAR-1/clone | jq . || true
	@echo ""
	@echo "5. Тарифы"
	@curl -sf localhost:$(PORT)/api/v1/marketplaces | jq .
	@echo ""
	@echo "$(GREEN)✓ smoke test завершён$(RESET)"

# ============================================================
# Полный цикл
# ============================================================

.PHONY: ci
ci: deps fmt vet test-race ## Полный локальный прогон как в CI
	@echo "$(GREEN)✓ CI-прогон завершён$(RESET)"

.PHONY: clean
clean: ## Удалить ./bin и coverage-файлы
	@rm -rf bin coverage.out coverage.html
	@echo "$(GREEN)✓ очищено$(RESET)"

.PHONY: all
all: build test ## Собрать и прогнать тесты