.PHONY = build test

LOG_DIR=./logs
CHECK_DIR=go list ./... | grep -v /cmd/utilits

generate-api:
	go get -u github.com/swaggo/swag/cmd/swag
	swag init -g ./cmd/server/main.go -o docs

build: generate-api
	go build -v ./cmd/server

build-docker:
	docker build --no-cache --network host -f ./docker/builder.Dockerfile . --tag patreon

run:
	#sudo chown -R 5050:5050 ./pgadmin
	mkdir -p $(LOG_DIR)
	docker-compose up --build --no-deps

run-with-build: build-docker run

open-last-log:
	cat $(LOG_DIR)/`ls -t $(LOG_DIR) | head -1 `

clear-logs:
	rm -r $(LOG_DIR)/*.out

stop:
	docker-compose stop

rm-docker:
	docker rm -vf $$(docker ps -a -q) || true

run-coverage:
	go test -covermode=atomic -coverpkg=$(go list ./... | grep -v "/cmd/utilits") -coverprofile=cover ./... | grep -v /cmd/utilits
	cat cover | fgrep -v "mock" | fgrep -v "testing.go" | fgrep -v "docs"  | fgrep -v "config" | fgrep -v "main" > cover2
	go tool cover -func=cover2

parse-last-log:
	go build -o logger.out -v ./cmd/utilits
	./logger.out -level=warn

test:
	go test -v -race ./... | grep -v /cmd/utilits
