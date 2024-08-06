TEST_FLAGS=""

build:
	@go build -o targets/ .

test:
	@go test $(TEST_FLAGS) ./...

exhaustive:
	@go run github.com/nishanths/exhaustive/cmd/exhaustive@latest ./...

sumtype:
	@go run github.com/BurntSushi/go-sumtype ./...

.PHONY: test

