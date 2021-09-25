.PHONY = build test

build:
	go build -v ./cmd/server

build-docker:
	docker build --no-cache --network host -f ./docker/builder.Dockerfile . --tag patreon

run:
	sudo chown -R 5050:5050 ./pgadmin
	docker-compose up --build --no-deps

stop:
	docker-compose down

rm-docker:
	docker rm -vf $$(docker ps -a -q) || true

test:
	go test -v -race ./...
