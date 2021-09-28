FROM golang:1.17.1

WORKDIR /app

COPY . .

EXPOSE 80 8080

RUN make build
RUN mkdir patreon-secrt

CMD ls
