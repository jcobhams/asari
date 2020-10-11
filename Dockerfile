FROM golang:1.14-stretch

ADD . /asari
WORKDIR /asari

CMD go test -v -race -coverprofile=cover.out ./...