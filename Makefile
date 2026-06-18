# ─────────────────────────────────────────────────────────────
# Makefile для avito_bonus_point_service
# ─────────────────────────────────────────────────────────────

DB_DSN := postgres://bonus:bonus@localhost:5432/bonus_ledger?sslmode=disable

# ── База данных ───────────────────────────────────────────────

## Поднять только Postgres в фоне
db-up:
	docker compose up -d postgres

## Остановить и удалить все контейнеры (данные сохраняются в volume)
db-down:
	docker compose down

## Полный сброс: удалить контейнеры И volume с данными
db-reset:
	docker compose down -v

# ── Тесты ────────────────────────────────────────────────────

## Запустить ВСЕ интеграционные тесты
test:
	TEST_DATABASE_URL="$(DB_DSN)" go test ./internal/api/... -v -count=1

## Запустить только тесты из Issue #8 (конкурентность + идемпотентность)
test-concurrent:
	TEST_DATABASE_URL="$(DB_DSN)" go test ./internal/api/... -v -count=1 \
		-run "TestConcurrent|TestConcurrentSameKey|TestConcurrentHolds"

## Запустить конкретный тест по имени (make test-one NAME=TestConcurrentDebitRaceCondition)
test-one:
	TEST_DATABASE_URL="$(DB_DSN)" go test ./internal/api/... -v -count=1 -run "$(NAME)"

# ── Полный цикл (самый простой способ) ───────────────────────

## Поднять БД → подождать → запустить тесты Issue #8 → остановить БД
ci:
	docker compose up -d postgres
	@echo "⏳ Жду готовности Postgres..."
	@until docker compose exec -T postgres pg_isready -U bonus -q; do sleep 1; done
	@echo "✅ Postgres готов"
	TEST_DATABASE_URL="$(DB_DSN)" go test ./internal/api/... -v -count=1 \
		-run "TestConcurrent" ; \
	EXIT_CODE=$$? ; \
	docker compose down ; \
	exit $$EXIT_CODE

.PHONY: db-up db-down db-reset test test-concurrent test-one ci
