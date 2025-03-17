FROM golang:1.24-alpine

RUN apk update && apk upgrade

WORKDIR /app

COPY . .

RUN go build -o main

CMD ["./main"]
