.PHONY: all
all: generate test

.PHONY: clean
clean:
	rm -rf gobin build

gobin:
	mkdir -p gobin

gobin/stringer: gobin go.sum go.mod
	GOBIN=$(CURDIR)/gobin go install golang.org/x/tools/cmd/stringer@latest

.PHONY: test
test:
	go test -mod vendor ./...

.PHONY: testcover
testcover:
	go test -mod vendor -race -coverprofile=coverage.txt -covermode=atomic ./...

.PHONY: generate
generate: gobin/stringer
	PATH=$(CURDIR)/gobin:$(PATH) go generate -mod vendor ./...

.PHONY: bench
bench:
	cd ./examples/bench/ && go test -bench=. -benchtime 5s -benchmem -run notest
# go test -v ./examples/bench/ --benchtime 10000x --bench ./examples/bench/ -benchmem

.PHONY: build
build:
	go build -o /dev/null ./...

.PHONY: tidy
tidy:
	go mod tidy
	go mod vendor
	rm -rf build gobin
