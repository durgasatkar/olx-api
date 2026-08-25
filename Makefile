.PHONY: build run

build:
	@go build -o bin/api ./cmd/api

run:
	@./bin/api

dev:
	@go run cmd/api/main.go
