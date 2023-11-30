.PHONY: build run test migrate-up migrate-down help sqlc mock

# Set your environment variables here or override them when calling `make <target>`.
ifneq ($(wildcard .env),)
	include .env
else
	echo "No .env file found. Please consider creating one."
endif

help: ## Show this help.
	@sed -ne '/@sed/!s/## //p' $(MAKEFILE_LIST) | column -tl 2

build: clean ## Build the application.
	go build ./...

clean: ## Clean the application.
	rm go-rssaggregator

db-up: ## Start the database and pgadmin in docker.
	docker-compose up

db-down: ## Stop the database and pgadmin.
	docker-compose down && docker system prune -f

run: ## Run the application.
	go run ./...

test: ## Test the application.
	go test ./...

migrate-up: ## Apply all up migrations.
	goose -dir sql/schema postgres $(DB_URL) up

migrate-down: ## Apply all down migrations.
	goose -dir sql/schema postgres $(DB_URL) down

sqlc: ## Generate the database code.
	sqlc generate

mock: ## Generate a store mock.
	mockgen -package mockdb -destination internal/mock/db.go github.com/jbdoumenjou/go-rssaggregator/internal/database Querier
