FROM golang:1.20.4-alpine as builder

WORKDIR /usr/app/game

RUN apk update && apk upgrade && \
    apk add --no-cache git

COPY ./go.* ./

RUN go mod download

COPY . .

RUN go build -o .bin/main cmd/main.go

FROM alpine

WORKDIR /usr/app/game

RUN apk update && apk upgrade && \
    apk add --no-cache git bash

COPY --from=builder /usr/app/game/.bin/main ./.bin/main
COPY --from=builder /usr/app/game/configs ./configs
COPY --from=builder /usr/app/game/migrations ./migrations

CMD [".bin/main"]