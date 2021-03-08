SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

export PATH := ./bin:$(PATH)
export GO111MODULE := on
export GOPROXY = https://proxy.golang.org,direct

# Install all the build and lint dependencies
setup:
	go mod download
	go generate -v ./...
.PHONY: setup

# Run all the tests
test:
	LC_ALL=C go test $(TEST_OPTIONS) -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=5m
.PHONY: test

# Run all the tests and opens the coverage report
cover: test
	go tool cover -html=coverage.txt
.PHONY: cover

# Clean up the Go module files
tidy:
	go mod tidy
.PHONY: tidy

# Upgrade non-test dependencies for the Go module
bump-deps: check
	go get -u ./...
.PHONY: bump-deps

# Upgrade all dependencies for the Go module
bump-deps-full: check
	go get -t -u ./...
.PHONY: bump-deps-full

# Inspect the source code for potential issues
vet:
	go vet ./...
.PHONY: vet

# Format source code to meet the language standards
fmt:
	go fmt ./...
.PHONY: fmt

# Run all the linters
lint:
	golangci-lint run ./...
	misspell -error **/*
.PHONY: lint

# Run all the tests and code checks
ci: build test vet lint
.PHONY: ci

# Build a beta version of updog
build:
	go build
.PHONY: build

# Show to-do items per file.
todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

.DEFAULT_GOAL := build