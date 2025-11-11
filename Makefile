SERVER_ADDRESS ?= localhost
SERVER_PORT ?= 8080
APP_NAME = shipping-optimizer
DOCKER_TAG = shipping-optimizer:latest

.PHONY: build run test docker-build docker-up docker-down fmt vet

build:
	go build -o ./bin/$(APP_NAME) ./cmd/api

run:
	@echo "Running app on $(SERVER_ADDRESS):$(SERVER_PORT)"
	SERVER_ADDRESS=$(SERVER_ADDRESS) SERVER_PORT=$(SERVER_PORT) go run ./cmd/api

test:
	go test -failfast -coverprofile=test.cover ./internal/... -v
	go tool cover -func=test.cover

fmt:
	go fmt ./...

vet:
	go vet ./...

docker-build:
	docker build -t $(DOCKER_TAG) -f misc/docker/dockerfile .

docker-up:
	docker compose -f ./misc/docker/docker-compose.yml up

docker-down:
	docker compose -f ./misc/docker/docker-compose.yml down
