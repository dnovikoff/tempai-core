
all: generate test build

clean:
	rm -rf gobin build

gobin:
	mkdir -p gobin

gobin/stringer: gobin
	go build -o ./gobin/stringer ./vendor/golang.org/x/tools/cmd/stringer

test:
	go test ./...

generate: gobin/stringer
	PATH=$(CURDIR)/gobin:$(PATH) go generate ./...

build:
	mkdir -p build
	GOBIN=$(shell pwd)/build go install ./examples/performance/...

