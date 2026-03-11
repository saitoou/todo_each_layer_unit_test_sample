SHELL := /bin/bash

.PHONY: help setup oapi-codegen generate statickcheck lint oapi-codegen

help: ## display thos help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

setup: ##install develop tools
	@go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.2.0
	@go install github.com/google/wire/cmd/wire@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.2.1
	@go install go.uber.org/mock/mockgen@latest

generate: ## generate to resolve depenedency injection from wire.go
	go generate ./...

lint: ## go run linter
	staticcheck ./...

docker-run: ## docker build and run
	docker build --no-cache -t todo-api -f ./docker/Dockerfile .
	docker run -it -rm -p 8080:8080 todo-api

docker-compose: ## docker-compose up
	docker compose -f docker/docker-compose.yaml down && docker compose -f docker/docker-compose.yaml up

oapi-codegen: ## genarate openapi models & servers
	oapi-codegen -config api/config/servers.config.yaml api/openapi.yaml
	oapi-codegen -config api/config/models.config.yaml api/openapi.yaml
