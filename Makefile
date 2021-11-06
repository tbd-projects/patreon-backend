.PHONY = build test

LOG_DIR=./logs
CHECK_DIR=go list ./... | grep -v /cmd/utilits

stop-redis:
	systemctl stop redis
stop-postgres:
	systemctl stop postgresql

run-posgres-redis:
	systemctl start redis
	systemctl start postgresql

watch-postgress-log:
	docker attach 2021_2_pyaterochka_patreon-bd_1

generate-api:
	go get -u github.com/swaggo/swag/cmd/swag
	swag init -g ./cmd/server/main.go -o docs

build: generate-api
	mkdir -p ./patreon-secrt
	go build -o server.out -v ./cmd/server

build-docker-local:
	docker build --no-cache --network host -f ./docker/builder.Dockerfile . --tag patreon

build-docker-server:
	docker build --build-arg RUN_HTTPS=-run-https --no-cache --network host -f ./docker/builder.Dockerfile . --tag patreon


run:
	#sudo chown -R 5050:5050 ./pgadmin
	mkdir -p $(LOG_DIR)
	docker-compose up --build --no-deps

run-with-build-local: build-docker-local run

run-with-build-server: build-docker-server run

open-last-log:
	cat $(LOG_DIR)/`ls -t $(LOG_DIR) | head -1 `

clear-logs:
	rm -r $(LOG_DIR)/*.log

stop:
	docker-compose stop

rm-docker:
	docker rm -vf $$(docker ps -a -q) || true

run-coverage:
	go test -covermode=atomic -coverpkg=./internal/... -coverprofile=cover ./internal/...
	cat cover | fgrep -v "mock" | fgrep -v "testing.go" | fgrep -v "docs"  | fgrep -v "config" | fgrep -v "main.go" > cover2
	go tool cover -func=cover2

build-utils:
	go build -o utils.out -v ./cmd/utilits

parse-last-log: build-utils
	./utils.out -search-url=${search_url}

gen-mock:
	go generate ./...

test:
	go test -v -race ./internal/...
