FROM golang:1.17.1

WORKDIR /app

COPY . .

EXPOSE 8080 443 80

RUN make build-sessions

CMD ./sessions.out