DATE ?= $(shell date +%FT%T%z)
BUILD_VERSION ?= '1.0.0'

BINARY=build/main

.PHONY: all
all: build test lint

.PHONE: configure
configure:
	go mod download && go mod tidy

.PHONY: build
build:
	 go build -o $(BINARY) ./cmd/main.go

.PHONY: test
test:
	 go test  ./internal/...

.PHONY: test.integration
test.integration:
	go test -v ./test/integration/...

.PHONY: lint
lint:
	golint ./...

.PHONY: clean
clean:
	rm -f $(BINARY)