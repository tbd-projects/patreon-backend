.PHONY = build test

generate-api:
	go get -u github.com/swaggo/swag/cmd/swag
	swag init -g ./cmd/server/main.go -o docs

build: generate-api
	go build -v ./cmd/server

build-docker:
	docker build --no-cache --network host -f ./docker/builder.Dockerfile . --tag patreon

run:
	#sudo chown -R 5050:5050 ./pgadmin
	mkdir -p ./logs
	docker-compose up --build --no-deps

run-with-build: build-docker run

stop:
	docker-compose stop

rm-docker:
	docker rm -vf $$(docker ps -a -q) || true

run-coverage:
	go test -covermode=atomic -coverpkg=./... -coverprofile=cover ./...
	cat cover | fgrep -v "mock" | fgrep -v "testing.go" | fgrep -v "docs"  | fgrep -v "config" | fgrep -v "main" > cover2
	go tool cover -func=cover2

test:
	go test -v -race ./...
