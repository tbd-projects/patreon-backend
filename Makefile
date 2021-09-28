.PHONY = build test

build:
	go build -v ./cmd/server

build-docker:
	docker build --no-cache --network host -f ./docker/builder.Dockerfile . --tag patreon

make-dir-cert:
	if ![ -f ./bin/status ]; then \
        mkdir ./patreon-secrt; \
	fi

run: make-dir-cert
	#sudo chown -R 5050:5050 ./pgadmin
	docker-compose up --build --no-deps

stop:
	docker-compose stop

rm-docker:
	docker rm -vf $$(docker ps -a -q) || true

test:
	go test -v -race ./...
