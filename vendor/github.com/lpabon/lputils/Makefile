all: build verify

.PHONY: build
build:
	go build $(shell go list ./... | grep -v vendor)

.PHONY: verify
verify: test

.PHONY: test
test:
	go test ./...