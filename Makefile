GOLANGCI_LINT = github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
TEST_FLAGS=""

build:
	@go build -o targets/ .

test:
	@go test $(TEST_FLAGS) ./...

lint.golangci-lint:
	go run $(GOLANGCI_LINT) run ./...

lint: lint.golangci-lint

