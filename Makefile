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
	go run ./cmd/migrate # optional if we add a migrate command later

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v 