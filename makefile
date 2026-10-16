DELIVERY_URL ?= http://localhost:8081
PORT         ?= 8080

GREEN  := \033[0;32m
YELLOW := \033[0;33m
RESET  := \033[0m

.PHONY: help
help: ## Список команд
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-22s$(RESET) %s\n", $$1, $$2}'

# ---------- всё в Docker ----------

.PHONY: up
up: ## Поднять всё в Docker (postgres + delivery + marketplace)
	docker compose up --build -d
	@echo "$(GREEN)✓ запущено: marketplace :8080, delivery :8081$(RESET)"

.PHONY: down
down: ## Остановить всё
	docker compose down

.PHONY: clean
clean: ## Остановить + удалить volume (сброс БД)
	docker compose down -v

.PHONY: logs
logs: ## Логи всех сервисов
	docker compose logs -f

.PHONY: ps
ps: ## Список контейнеров
	docker compose ps

# ---------- локально ----------

.PHONY: delivery
delivery: ## Запустить delivery локально
	cd delivery && go run ./cmd/api

.PHONY: marketplace
marketplace: ## Запустить marketplace локально (in-memory)
	cd marketplace && STORAGE=memory DELIVERY_URL=$(DELIVERY_URL) PORT=$(PORT) go run ./cmd/api

.PHONY: marketplace-postgres
marketplace-postgres: ## Запустить marketplace локально с Postgres
	cd marketplace && STORAGE=postgres \
		DATABASE_URL='postgres://marketplace:marketplace@localhost:5432/marketplace?sslmode=disable' \
		DELIVERY_URL=$(DELIVERY_URL) PORT=$(PORT) go run ./cmd/api

# ---------- тесты и сборка ----------

.PHONY: build
build: ## Собрать оба сервиса
	cd marketplace && go build ./...
	cd delivery && go build ./...
	@echo "$(GREEN)✓ собрано$(RESET)"

.PHONY: test
test: ## Прогнать тесты обоих сервисов
	cd marketplace && go test ./...
	cd delivery && go test ./...

.PHONY: tidy
tidy: ## go mod tidy для обоих
	cd marketplace && go mod tidy
	cd delivery && go mod tidy

.PHONY: fmt
fmt: ## Форматирование
	cd marketplace && go fmt ./...
	cd delivery && go fmt ./...

# ---------- Postgres отдельно ----------

.PHONY: postgres-up
postgres-up: ## Поднять только Postgres
	docker compose up -d postgres

.PHONY: postgres-down
postgres-down: ## Остановить Postgres
	docker compose stop postgres

.PHONY: psql
psql: ## Открыть psql
	docker exec -it marketplace-postgres psql -U marketplace -d marketplace

.PHONY: db-listings
db-listings: ## Показать объявления в БД
	docker exec -it marketplace-postgres psql -U marketplace -d marketplace \
		-c "SELECT id, title, status, delivery_id FROM listings;"

# ---------- smoke ----------

.PHONY: smoke
smoke: ## Прогнать smoke-сценарий
	@echo "$(GREEN)1. health$(RESET)"
	@curl -sf localhost:8080/healthz | jq
	@curl -sf localhost:8081/healthz | jq
	@echo "$(GREEN)2. создать объявление$(RESET)"
	@curl -sf -X POST localhost:8080/api/v1/listings \
		-H "Content-Type: application/json" \
		-d '{"car_id":"CAR-1","brand":"Tesla","model":"Model 3","year":2023,"mileage":15000,"price":5000000,"seller_id":"D-1","seller_name":"AutoDealer","seller_type":"dealer"}' | jq
	@echo "$(GREEN)3. опубликовать$(RESET)"
	@curl -sf -X POST localhost:8080/api/v1/listings/L-D-1-CAR-1/publish | jq
	@echo "$(GREEN)4. продать (marketplace → delivery)$(RESET)"
	@curl -sf -X POST localhost:8080/api/v1/listings/L-D-1-CAR-1/sell | jq
	@echo "$(GREEN)5. доставки$(RESET)"
	@curl -sf localhost:8081/api/v1/deliveries | jq
	@echo "$(GREEN)6. запустить доставку$(RESET)"
	@curl -sf -X POST localhost:8081/api/v1/deliveries/DEL-L-D-1-CAR-1/start | jq
	@echo "$(GREEN)7. завершить$(RESET)"
	@curl -sf -X POST localhost:8081/api/v1/deliveries/DEL-L-D-1-CAR-1/complete | jq
	@echo "$(GREEN)8. отчёты (Visitor)$(RESET)"
	@curl -sf localhost:8081/api/v1/reports | jq
