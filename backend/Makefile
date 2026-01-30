.PHONY: run docs

run:
	go run cmd/app/main.go

docs:
	swag init -g ./internal/delivery/http/v1/handler.go -o ./docs	