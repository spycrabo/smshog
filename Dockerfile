# smshog Dockerfile

FROM golang:1.19-alpine AS builder

LABEL maintainer="Vladimir Tretiakov"

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -o smshog main.go

FROM golang:1.19-alpine

WORKDIR /build

COPY --from=builder /build/smshog /build/smshog

CMD ["./smshog"]