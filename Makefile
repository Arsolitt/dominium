build:
	@go build -o bin/dominium cmd/main.go

run: build
	@./bin/dominium | jq

test:
	@go test -v ./...
