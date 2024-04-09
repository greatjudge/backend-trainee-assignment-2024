ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=test password=test dbname=test host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/internal
MIGRATION_FOLDER=$(INTERNAL_PKG_PATH)/db/migrations

.PHONY: migration-create
migration-create:
	goose -dir "$(MIGRATION_FOLDER)" -s create "$(name)" sql

.PHONY: test-migration-up
test-migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

.PHONY: test-migration-down
test-migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

.PHONY: test-db-up
test-db-up:
	docker compose up -d

.PHONY: test-db-down
test-db-down:
	docker compose down -v

.PHONY: run-app
run-app:
	make test-migration-up
	go run cmd/main.go

.test-integration:
	$(info Running integration tests...)
	go test -v -coverprofile=cover.out ./tests/... -cover
	go tool cover -html=cover.out -o cover.html
test-integr: .test-integration ##запуск всех интеграционных тестов