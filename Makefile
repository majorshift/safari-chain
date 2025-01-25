build:
	@go build -o bin/safarichain

run: build
	@./bin/safarichain

test:
	@go test -v ./...
