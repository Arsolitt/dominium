ifneq (,$(wildcard ./.env))
	include .env
	export
endif

build:
	@go build -o bin/dominium main.go

run: build
	@./bin/dominium | jq

test:
	@go test -v ./...
