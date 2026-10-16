.PHONY: init run build fmt test vet check

init:
	go mod download

run:
	go run ./cmd/guardian

build:
	go build -o bin/guardian ./cmd/guardian

fmt:
	gofmt -w cmd internal

test:
	go test ./...

vet:
	go vet ./...

check: test vet
