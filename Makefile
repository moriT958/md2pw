.PHONY: build test
build: main.go
	go build -o md2puki .

test:
	go test -v ./...
