FROM golang:1.17.1

WORKDIR /app

COPY . .

EXPOSE 80 8081

RUN make build

CMD ./server