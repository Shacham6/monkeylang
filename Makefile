GOLANGCI_LINT = github.com/golangci/golangci-lint/cmd/golangci-lint@v1.61.0
TEST_FLAGS=""

build:
	@go build -o targets/ .

test:
	@go test $(TEST_FLAGS) ./...

exhaustive:
	@go run github.com/nishanths/exhaustive/cmd/exhaustive@latest ./...

sumtype:
	@go run github.com/BurntSushi/go-sumtype ./...

lint.golangci-lint:
	go run $(GOLANGCI_LINT) run ./...

lint: lint.golangci-lint

