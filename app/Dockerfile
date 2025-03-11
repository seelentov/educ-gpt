FROM golang:1.24-alpine

RUN apk update && apk upgrade

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["go","run","."]

