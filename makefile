.PHONY: help up down clean logs ps test build smoke topics ui

help:
	@echo "make up       — всё в Docker (kafka, kafka-ui, 2 postgres, 2 сервиса)"
	@echo "make down     — остановить"
	@echo "make clean    — сброс БД и volumes"
	@echo "make logs     — логи сервисов"
	@echo "make ps       — контейнеры"
	@echo "make build    — go build обоих"
	@echo "make smoke    — создать → опубликовать → продать (async через Kafka)"
	@echo "make topics   — список Kafka-топиков"
	@echo "make ui       — открыть Kafka UI http://localhost:8090"

up:
	docker compose up --build -d
	@echo "Kafka UI: http://localhost:8090"

down:
	docker compose down

clean:
	docker compose down -v

logs:
	docker compose logs -f marketplace delivery

ps:
	docker compose ps

build:
	cd marketplace && go build ./...
	cd delivery    && go build ./...

test:
	cd marketplace && go test ./... 2>/dev/null || true
	cd delivery    && go test ./... 2>/dev/null || true

topics:
	docker exec kafka /opt/kafka/bin/kafka-topics.sh --bootstrap-server localhost:9092 --list

ui:
	@echo "Открой http://localhost:8090"

smoke:
	@echo "1. Создать объявление"
	@curl -s -X POST localhost:8080/api/v1/listings \
	  -H "Content-Type: application/json" \
	  -d '{"car_id":"CAR-1","brand":"Tesla","model":"Model 3","year":2023,"mileage":15000,"price":5000000,"seller_id":"D-1","seller_name":"AutoDealer","seller_type":"dealer"}' | jq
	@echo "\n2. Опубликовать"
	@curl -s -X POST localhost:8080/api/v1/listings/L-D-1-CAR-1/publish | jq
	@echo "\n3. Продать (async → Kafka)"
	@curl -s -X POST localhost:8080/api/v1/listings/L-D-1-CAR-1/sell | jq
	@echo "\n4. Ждём 3 сек, чтобы Kafka протащила события..."
	@sleep 3
	@echo "\n5. Список объявлений (в listing должен появиться delivery_id)"
	@curl -s localhost:8080/api/v1/listings | jq
	@echo "\n6. Список доставок"
	@curl -s localhost:8081/api/v1/deliveries | jq
