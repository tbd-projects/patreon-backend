FROM golang:1.17.1

WORKDIR /app

COPY . .

RUN apt-get update
RUN apt-get install jq -y

EXPOSE 443 80 8080 9001

RUN make build

ARG RUN_HTTPS

ENV HTTPS=$RUN_HTTPS

RUN chmod +x ./wait

CMD ./wait && ./server.out -server-run $HTTPS