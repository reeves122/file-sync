FROM golang:1.17-alpine AS builder

COPY . /usr/src/app
WORKDIR /usr/src/app

RUN go get -v -t -d ./...
RUN go build -o app .

FROM alpine:latest

RUN apk update
RUN apk add git --no-cache

WORKDIR /
COPY --from=builder /usr/src/app/app .

CMD ["/app"]