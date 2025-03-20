FROM golang:1.24-alpine

RUN apk update && apk upgrade && \
    apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

ENV CGO_ENABLED=1
RUN go build -o main

RUN go test test/unit/*

# RUN go test test/integ/*

CMD ["./main"]
