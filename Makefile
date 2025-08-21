.PHONY: run build tidy gqlgen wire docker-up docker-down migrate

run:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

clean:
	rm -rf bin

tidy:
	go mod tidy

wire:
	go run github.com/google/wire/cmd/wire@v0.6.0 ./...

gqlgen:
	go run github.com/99designs/gqlgen generate

migrate:
	go run ./cmd/migrate -action=up

migrate-down:
	go run ./cmd/migrate -action=down

migrate-status:
	go run ./cmd/migrate -action=status

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v

docker-debug:
	docker compose -f docker-compose.debug.yml up --build -d

docker-debug-down:
	docker compose -f docker-compose.debug.yml down -v

docker-dev:
	docker compose -f docker-compose.dev.yml up --build -d

docker-dev-down:
	docker compose -f docker-compose.dev.yml down -v 