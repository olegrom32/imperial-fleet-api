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

Please first run `docker-compose up db` and wait until the db is fully initialized.

After that please run `docker-compose up api`. The service will be exposed on `:8080` port.

## Unit tests

Only one handler is covered as an example (`/internal/applciation/handler/updatespaceship/handler_test.go`).

## Authentication

A simple basic auth is added, here are the authorized users:
- user: user1, password: test1
- user: user2, password: test2
