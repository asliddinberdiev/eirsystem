.PHONY: run docs migrate-create migrate-up migrate-down migrate-status migrate-reset

run:
	go run cmd/app/main.go

docs:
	swag init -g ./internal/delivery/http/v1/handler.go -o ./docs
	./scripts/clean_docs.sh	

# Migration commands
# Usage: make migrate-create name=create_users_table
migrate-create:
	@if [ -z "$(name)" ]; then echo "Error: name is required. Usage: make migrate-create name=migration_name"; exit 1; fi
	goose -dir ./migrations create $(name) sql

migrate-up:
	@. ./.env && goose -dir ./migrations postgres "host=$$POSTGRES_HOST user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DBNAME port=$$POSTGRES_PORT sslmode=$$POSTGRES_SSLMODE" up

migrate-down:
	@. ./.env && goose -dir ./migrations postgres "host=$$POSTGRES_HOST user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DBNAME port=$$POSTGRES_PORT sslmode=$$POSTGRES_SSLMODE" down

migrate-status:
	@. ./.env && goose -dir ./migrations postgres "host=$$POSTGRES_HOST user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DBNAME port=$$POSTGRES_PORT sslmode=$$POSTGRES_SSLMODE" status

migrate-reset:
	@. ./.env && goose -dir ./migrations postgres "host=$$POSTGRES_HOST user=$$POSTGRES_USER password=$$POSTGRES_PASSWORD dbname=$$POSTGRES_DBNAME port=$$POSTGRES_PORT sslmode=$$POSTGRES_SSLMODE" reset