# syntax=docker/dockerfile:1
FROM golang:1.24.3-alpine3.21 AS builder

WORKDIR /go/delivery

COPY . .

RUN go build -o GoBillIt ./cmd/GoBillIt/

FROM alpine:3.21.3
WORKDIR /GoBillIt
COPY --from=builder /go/delivery/GoBillIt ./
CMD ["./GoBillIt"]
