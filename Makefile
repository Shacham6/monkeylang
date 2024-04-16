build:
	@go build -o targets/ .

test:
	@go test ./...
.PHONY: test

