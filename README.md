# imperial-fleet-api
Tech assignment

## Dir structure

This repo follows the https://github.com/golang-standards/project-layout

## Internal packages

Hexagonal/DDD approach (simplified) is attempted here.

- `/internal/application` holds the application/business logic/domain
- `/internal/infra` holds the infrastructure adapters (repos)
- `/internal/server` contains REST/HTTP protocol adapters

## Makefile

A couple of handy commands are available via `make`

- `make generate` runs all code generators
- `make test` runs all tests
- `make migrate` runs all migrations

## Starting locally

Just run `docker-compose up` and the API will be exposed on `:8080` port after everything is initialized.

## Unit tests

Only one handler is covered as an example (`/internal/applciation/handler/updatespaceship/handler_test.go`).

## Authentication

A simple basic auth is added, here are the authorized users:
- user: user1, password: test1
- user: user2, password: test2

## API documentation (Swagger)

After docker compose has started please navigate to http://localhost:8088/ for Swagger API documentation

## TODOs

The `armament` entity and the `spaceship->armament` relation is never populated - did not have enough time to do that.

No logger

Creation does not check duplicates

Unit tests coverage must be improved

Functional tests must be added in addition to unit tests

Linter must be added

Github Actions to run tests and linters before any merge to `main` must be added
