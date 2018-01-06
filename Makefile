
all: generate test build

test:
	go test ./...

generate:
	go generate ./...

build:
	rm -rf build
	mkdir -p build
	GOBIN=$(shell pwd)/build go install ./examples/...
