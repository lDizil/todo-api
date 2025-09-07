FROM golang:latest AS builder 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o todoapi main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todoapi .

CMD ["./todoapi"]