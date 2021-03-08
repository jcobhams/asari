FROM golang:1.14-alpine
RUN apk add build-base

ADD . /asari
WORKDIR /asari

CMD go test -v -race -coverprofile=cover.out ./...