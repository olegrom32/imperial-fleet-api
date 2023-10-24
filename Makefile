bin/: ; mkdir -p $@

.PHONY: migrate
migrate: bin/migrate
	bin/migrate -source file://./migrations -database "mysql://root:@tcp(localhost:3306)/api" up

bin/migrate: | bin/
	GOBIN="$(realpath $(dir $@))" go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.16.2

.PHONY: generate
generate: bin/mockgen
	go generate ./...

bin/mockgen: | bin/
	GOBIN="$(realpath $(dir $@))" go install go.uber.org/mock/mockgen@v0.2.0

.PHONY: test
test:
	go test -v -race ./internal/...
