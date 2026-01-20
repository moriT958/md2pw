.PHONY: cli test lint
cli:
	go build -o md2pw ./cmd/cli

test:
	go test -v ./...

lint:
	golangci-lint run ./...
