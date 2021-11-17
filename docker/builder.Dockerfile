FROM golang:1.17.1

WORKDIR /app

COPY . .

RUN apt-get update
RUN apt-get install jq -y

EXPOSE 443 80 8080

RUN make build

ARG RUN_HTTPS

ENV HTTPS=$RUN_HTTPS

CMD ./server.out -server-run $HTTPS