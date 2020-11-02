
all: generate test build

clean:
	rm -rf gobin build

gobin:
	mkdir -p gobin

gobin/stringer: gobin
	go build -mod vendor -o ./gobin/stringer ./vendor/golang.org/x/tools/cmd/stringer

test:
	go test -mod vendor ./...

testcover:
	go test -mod vendor -race -coverprofile=coverage.txt -covermode=atomic ./...

generate: gobin/stringer
	PATH=$(CURDIR)/gobin:$(PATH) go generate -mod vendor ./...

build:
	mkdir -p build
	GOBIN=$(shell pwd)/build go install -mod vendor ./examples/performance/...

