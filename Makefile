ifeq ($(POSTGRES_SETUP_TEST),)
	POSTGRES_SETUP_TEST := user=test password=test dbname=test host=localhost port=5432 sslmode=disable
endif

INTERNAL_PKG_PATH=$(CURDIR)/app/internal
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

.PHONY: run-app
run-app:
	docker compose --env-file .env up

.PHONY: test-app-up
test-app-up:
	docker compose --env-file .env up -d

.PHONY: test-app-down
test-app-down:
	docker compose --env-file .env down -v

.test-integration:
	$(info Running integration tests...)
	cd app && go test -v -coverprofile=cover.out ./tests/... -cover
	cd app && go tool cover -html=cover.out -o cover.html
test-integr: .test-integration ##запуск всех интеграционных тестов