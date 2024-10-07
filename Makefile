GOLANGCI_LINT = github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
TEST_FLAGS=""

export GOBIN ?= $(shell pwd)/.bin/
export PATH := ${GOBIN}:${PATH}

all:
	$(MAKE) -k lint test

build:
	@go build -o targets/ .

test:
	@go test $(TEST_FLAGS) ./...

lint.install.golangci-lint:
	go install $(GOLANGCI_LINT)

lint.install: lint.install.golangci-lint

lint.golangci-lint:
	golangci-lint run ./...

lint: lint.golangci-lint

