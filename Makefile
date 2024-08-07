# Load tasks.
-include tasks/Makefile.*

.DEFAULT_GOAL := help: ## Prints help for targets with comments

.PHONY: help
## Display this help message
help:
	@scripts/print_make_help.sh $(shell realpath $(MAKEFILE_LIST))

.PHONY: mock
## Generate mock interfaces for testing
mock: mock/user_service mock/user_repository

.PHONY: build
## Build an optimized Docker image. Alias for docker/build.
build: docker/build

.PHONY: clean
## Remove all Make-generated artifacts.
clean: docker/clean generate/clean mock/clean

.PHONY: generate
## Generate development dependencies.
generate: generate/queries generate/data_mount_fixtures

.PHONY: psql
## Run psql against the database. Alias for docker/exec/psql.
psql: docker/exec/psql

.PHONY: postgres
## Run postgres db without other affiliated compose containers
postgres: docker/up/postgres

.PHONY: up
## Run the app interactively. Alias for docker/up.
up: docker/up

.PHONY: down
## Stop the app. Alias for docker/down.
down: docker/down

.PHONY: test
## Run the test suite. Alias for docker/test.
test: docker/test
