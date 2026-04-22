GO ?= go
BINARY ?= dewdb

.PHONY: fmt lint test build tidy

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run

test:
	$(GO) test ./...

build:
	$(GO) build -o bin/$(BINARY) ./cmd/dewdb

tidy:
	$(GO) mod tidy
