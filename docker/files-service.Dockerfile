FROM golang:1.17.1

WORKDIR /app

COPY . .

EXPOSE 8080 443 80 5002

RUN make build-files

CMD ./files.out