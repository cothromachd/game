FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

FROM alpine

WORKDIR /app

COPY configs/config.yaml /app/configs/config.yaml
COPY --from=builder /app/main /app/main

CMD ["/app/main"]
