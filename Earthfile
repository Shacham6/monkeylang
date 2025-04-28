VERSION 0.8
FROM golang:1.21.7

WORKDIR /app

build:
	COPY . .
	RUN go build -o targets/ .

test:
	COPY . .
	RUN go test ./...

