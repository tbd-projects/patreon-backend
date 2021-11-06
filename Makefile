.PHONY = build test

LOG_DIR=./logs
CHECK_DIR=go list ./... | grep -v /cmd/utilits
SQL_DIR=./scripts

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

$(serv):
	cat ./configs/migrate.config | jq '.database_server'

DATABASE_URL=$$(cat ./configs/migrate.config | jq '.database_server')
DATABASE_URL_LOCAL="$$(cat ./configs/migrate.config | jq '.database_local')"



some:
	echo ${DATABASE_URL}

migrate-up:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL} up

local-migrate-up:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL_LOCAL} up

migrate-down:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL} down

local-migrate-down:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL_LOCAL} down

migrate-up-one:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL} up 1

local-migrate-up-one:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL_LOCAL} up 1

migrate-down-one:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL} down 1

local-migrate-down-one:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL_LOCAL} down 1

migrate:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL} goto $(version)

local-migrate:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL_LOCAL} goto $(version)

force-migrate:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL} force $(version)

local-force-migrate:
	migrate -source file://${SQL_DIR} -database ${DATABASE_URL_LOCAL} force $(version)


