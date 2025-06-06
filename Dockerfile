FROM golang:1.21.7 as base
WORKDIR /app
COPY . .

FROM base as build
RUN ["go", "build", "-o", "targets/"]

FROM base as test
RUN ["go", "test", "./..."]

