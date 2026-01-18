.PHONY: build test
build:
	go build -o md2pw ./cmd/md2pw

test:
	go test -v ./...
