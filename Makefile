APP_NAME = sso
DOCKER_COMPOSE = docker-compose

.PHONY: build up down logs restart migrate-new migrate-up migrate-down

build:
	go build -o $(APP_NAME) ./cmd/app

up:
	$(DOCKER_COMPOSE) up -d --build

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f sso

restart: down up

migrate-new:
	migrate create -ext sql -dir migrations -seq init

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/sso?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/sso?sslmode=disable" down 1
