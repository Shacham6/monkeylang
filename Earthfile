VERSION 0.8
FROM golang:1.21.7

WORKDIR /app

build:
	COPY . .
	RUN make build

test:
	COPY . .
	RUN make test

