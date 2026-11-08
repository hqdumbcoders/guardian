FROM golang:1.27.2 AS builder

WORKDIR /dev/guardian

COPY go.mod ./

COPY . .
RUN go build -v -o ./app ./cmd/guardian

FROM alpine:3.23 AS runtime

WORKDIR /app

COPY --from=builder /dev/guardian/app ./app

CMD ["./app"]

