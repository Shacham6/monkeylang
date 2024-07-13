build:
	@go build -o targets/ .

test:
	@go test ./...

exhaustive:
	@go run github.com/nishanths/exhaustive/cmd/exhaustive@latest ./...

sumtype:
	@go run github.com/BurntSushi/go-sumtype ./...

.PHONY: test

