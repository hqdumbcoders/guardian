FROM golang:1.27.2 AS builder

WORKDIR /app

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 go build -v -o ./app ./cmd/guardian

FROM alpine:3.23 AS runtime

WORKDIR /app

COPY --from=builder /app/app ./app

CMD ["./app"]

