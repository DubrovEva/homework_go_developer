include make/lint.mk
include make/build.mk

ifndef VERBOSE
.SILENT:
endif

lint: cart-lint #loms-lint notifier-lint comments-lint

build: cart-build loms-build notifier-build comments-build

run-all:
	docker compose up --build

run-monitoring:
	docker compose up prometheus grafana jaeger --build

run-migrations:
	loms/bin/goose -dir "${MIGRATION_DIR}" postgres "postgresql://loms-user:loms-password@127.0.0.1:5432/loms_db?sslmode=disable" up

stop-all:
	docker compose down -v

test:
	go test ./cart/... -cover
