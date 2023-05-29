.PHONY: install test cover clean test-release

VERSION=$(shell git describe --always --dirty --tags)

install: test
	go install -ldflags="-X main.version=$(VERSION)" ./...

test: clean
	go test -v --timeout 10s -race ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf ./build
	mkdir ./build
	rm -f ./testdata10_*
	rm -f ./testdata50k_*

test-release:
	goreleaser release --snapshot --clean
