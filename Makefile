.PHONY: cli test
cli:
	go build -o md2pw ./cmd/cli

test:
	go test -v ./...
