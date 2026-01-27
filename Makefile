.PHONY: run backend frontend up down logs tidy

run:
	@make -j 2 backend frontend

up:
	cd infra && docker compose up -d

down:
	cd infra && docker compose down

backend:
	cd backend && go run cmd/main.go

frontend:
	cd frontend && bun run dev

logs:
	cd infra && docker compose logs -f

tidy:
	cd backend && go mod tidy