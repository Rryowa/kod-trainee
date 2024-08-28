include .env
DB_STRING="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"
MIGRATIONS_DIR="./migrations"

build:
	@go build -o bin/main cmd/main.go
	@chmod +x bin/main

run:
	@bin/main

all: build up run

up:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) up

down:
	@goose -dir $(MIGRATIONS_DIR) postgres $(DB_STRING) down