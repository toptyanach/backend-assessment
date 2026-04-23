.PHONY: test lint build docker

test:
	go test -v -cover ./...

lint:
	golangci-lint run ./...

build:
	go build -o bin/server ./cmd/server/main.go

docker:
	docker build -t backend-assessment:latest .