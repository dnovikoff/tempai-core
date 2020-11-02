.PHONY: all
all: generate test build

.PHONY: clean
clean:
	rm -rf gobin build

gobin:
	mkdir -p gobin

gobin/stringer: gobin
	go build -mod vendor -o ./gobin/stringer ./vendor/golang.org/x/tools/cmd/stringer

.PHONY: test
test:
	go test -mod vendor ./...

.PHONY: testcover
testcover:
	go test -mod vendor -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: generate
generate: gobin/stringer
	PATH=$(CURDIR)/gobin:$(PATH) go generate -mod vendor ./...

.PHONY: build
build:
	mkdir -p build
	GOBIN=$(shell pwd)/build go install -mod vendor ./examples/performance/...

