APP_NAME=books-api
PORT=8080

.PHONY: build run fmt

build:
	go mod tidy
	go build -o $(APP_NAME) ./cmd/api-service

run: build
	PORT=$(PORT) ./$(APP_NAME)

fmt:
	gofmt -w .
