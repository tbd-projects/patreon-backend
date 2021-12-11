.PHONY = build test

GRAFANA_DIR=./grafana
LOG_DIR=./logs
LOG_SESSION_DIR = ./logs-sessions
LOG_SESSION_DIR = ./logs-files
CHECK_DIR=go list ./... | grep -v /cmd/utilits
SQL_DIR=./scripts
MICROSERVICE_DIR=$(PWD)/internal/microservices

stop-redis:
	systemctl stop redis
stop-postgres:
	systemctl stop postgresql
run-posgres-redis:
	systemctl start redis
	systemctl start postgresql

generate-api:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/server/main.go -o docs

build: generate-api
	mkdir -p ./patreon-secrt
	go build -o server.out -v ./cmd/server
build-sessions:
	go build -o sessions.out -v ./cmd/sessions
build-files:
	go build -o files.out -v ./cmd/files

build-docker-server: # запуск обычного http servera
	docker build --no-cache --network host -f ./docker/builder.Dockerfile . --tag patreon
build-docker-server-https: # запуск https serverа
	docker build --build-arg RUN_HTTPS=-run-https --no-cache --network host -f ./docker/builder.Dockerfile . --tag patreon

build-docker-pg: # сборка образа базы
	docker build --no-cache --network host -f ./docker/postgresql.Dockerfile . --tag pg-14
build-docker-sessions: # сборка образа сервиса авторизаций
	docker build --no-cache --network host -f ./docker/session-service.Dockerfile . --tag session-service
build-docker-files: # сборка образа сервиса файлов
	docker build --no-cache --network host -f ./docker/files-service.Dockerfile . --tag files-service
build-docker-nginx: # сборка образа сервиса nginx
	docker build --no-cache --network host -f ./docker/nginx.Dockerfile . --tag nginx-ssl


run-init:
	#sudo chown -R 5050:5050 ./pgadmin
	mkdir -p $(LOG_DIR)
	mkdir -p $(GRAFANA_DIR)
	sudo chown -R 472:472 ./grafana

run-https: run-init # запустить https сервер
	docker-compose --env-file ./configs/run-https.env up --build --no-deps

run-http: run-init # запустить http сервер
	docker-compose --env-file ./configs/run-http.env up --build --no-deps

stop:  # остановить сервер
	docker-compose stop

# запустить http сервер с http nginx
run-with-build-http: build-docker-server build-docker-sessions build-docker-files build-docker-pg run-http

# запустить https сервер с http nginx
run-with-build-https: build-docker-server-https build-docker-sessions build-docker-files build-docker-pg run-http

# запустить http сервер с https nginx
run-with-build: build-docker-server build-docker-sessions build-docker-files build-docker-pg run-https


open-last-log:
	cat $(LOG_DIR)/`ls -t $(LOG_DIR) | head -1 `

watch-postgress-log:
	docker attach 2021_2_pyaterochka_patreon-bd_1

clear-logs:
	rm -rf $(LOG_DIR)/*.log
	rm -rf $(LOG_SESSION_DIR)/*.log
	rm -rf $(LOG_FILES_DIR)/*.log

rm-docker:
	docker rm -vf $$(docker ps -a -q) || true

run-coverage:
	go test -covermode=atomic -coverpkg=./internal/... -coverprofile=cover ./internal/...
	cat cover | fgrep -v "mock" | fgrep -v "testing.go" | fgrep -v "docs" | fgrep -v ".pb.go" | fgrep -v "config" |fgrep -v "patreon/internal/app/server/server.go" > cover2
	go tool cover -func=cover2

build-utils:
	go build -o utils.out -v ./cmd/utilits

parse-last-log: build-utils
	./utils.out -search-url=${search_url}

gen-mock:
	go generate ./...

gen-proto-sessions:
	protoc --proto_path=${MICROSERVICE_DIR}/auth/delivery/grpc/protobuf session.proto --go_out=plugins=grpc:${MICROSERVICE_DIR}/auth/delivery/grpc/protobuf
gen-proto-files:
	protoc --proto_path=${MICROSERVICE_DIR}/files/delivery/grpc/protobuf files.proto --go_out=plugins=grpc:${MICROSERVICE_DIR}/files/delivery/grpc/protobuf

test:
	go test -v -race ./internal/...


DATABASE_URL:=$(shell cat ./configs/migrate.config | jq '.database_server')
DATABASE_URL_LOCAL:=$(shell cat ./configs/migrate.config | jq '.database_local')

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


