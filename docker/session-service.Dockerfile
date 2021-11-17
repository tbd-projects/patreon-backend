FROM golang:1.17.1

WORKDIR /app

COPY . .

EXPOSE 5000

RUN make build-sessions

CMD ./sessions.out