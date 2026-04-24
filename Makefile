GO ?= go
BINARY ?= dewdb
CLI_BINARY ?= dewdb-cli

.PHONY: fmt lint test build build-cli tidy

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run

test:
	$(GO) test ./...

build:
	$(GO) build -o bin/$(BINARY) ./cmd/dewdb

build-cli:
	$(GO) build -o bin/$(CLI_BINARY) ./cmd/dewdb-cli

tidy:
	$(GO) mod tidy
