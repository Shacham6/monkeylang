FROM golang:1.21.7 as base
WORKDIR /app
COPY . .

FROM base AS build
RUN ["go", "build", "-o", "targets/"]

FROM base AS test
RUN ["go", "test", "./..."]

