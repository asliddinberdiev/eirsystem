.PHONY: run docs migrate-create migrate-up migrate-down migrate-status migrate-reset

run:
	go run cmd/app/main.go

docs:
	swag init -g internal/delivery/http/v1/handler.go --parseDependency --parseInternal --parseDepth 1 -o ./docs
	./scripts/clean_docs.sh

migrate-create:
	goose -s -dir ./migrations create $(name) sql

migrate-up:
	@. ./.env && goose -dir ./migrations postgres "$(dsn)" up

migrate-down:
	@. ./.env && goose -dir ./migrations postgres "$(dsn)" down

migrate-status:
	@. ./.env && goose -dir ./migrations postgres "$(dsn)" status

migrate-reset:
	@. ./.env && goose -dir ./migrations postgres "$(dsn)" reset