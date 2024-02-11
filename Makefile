.PHONY: build run test help sqlc mock

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

run: ## Run the application.
	go run ./...

test: ## Test the application.
	go test -v ./...

sqlc: ## Generate the database code.
	sqlc generate

mock: ## Generate a store mock.
	mockgen -package mockdb -destination internal/mock/db.go github.com/jbdoumenjou/go-rssaggregator/internal/database Querier
