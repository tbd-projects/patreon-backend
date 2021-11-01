FROM golang:1.17.1

WORKDIR /app

COPY . .

EXPOSE 443 80

RUN make build

ARG RUN_HTTPS

ENV HTTPS=$RUN_HTTPS

CMD ./server.out -server-run $HTTPS