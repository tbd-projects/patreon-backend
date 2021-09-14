.PHONY = build test

build:
	go build -v ./cmd/server

test:
	go test -v -race ./...
