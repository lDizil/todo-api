FROM golang:latest

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o todoapi .

CMD ["./todoapi"]