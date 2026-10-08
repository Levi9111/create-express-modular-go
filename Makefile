.PHONY: build test test-race test-coverage clean install build-all fmt vet check help

BINARY_NAME=cem
BUILD_DIR=bin
VERSION_PKG=github.com/Levi9111/create-express-modular-go/internal/version
GIT_COMMIT?=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS=-ldflags="-s -w -X $(VERSION_PKG).GitCommit=$(GIT_COMMIT) -X $(VERSION_PKG).BuildDate=$(BUILD_DATE)"

help:
	@echo "Available Make targets:"
	@echo "  make build          - Compile binary into $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "  make install        - Install binary globally to \$$GOPATH/bin"
	@echo "  make test           - Run all unit tests"
	@echo "  make test-race      - Run tests with race detector enabled"
	@echo "  make test-coverage  - Generate coverage profile and summary"
	@echo "  make fmt            - Run go fmt on all packages"
	@echo "  make vet            - Run go vet on all packages"
	@echo "  make check          - Run vet and test suite"
	@echo "  make build-all      - Cross-compile binaries for Linux, macOS, and Windows"
	@echo "  make clean          - Remove build artifacts and temporary binaries"

build:
	@mkdir -p $(BUILD_DIR)
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) main.go
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

install:
	go install $(LDFLAGS)
	@echo "Installed $(BINARY_NAME) to \$$GOPATH/bin"

test:
	go test -v ./...

test-race:
	go test -race -v ./...

test-coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

fmt:
	go fmt ./...

vet:
	go vet ./...

check: vet test

clean:
	rm -rf $(BUILD_DIR) dist coverage.out coverage.html

build-all:
	@mkdir -p $(BUILD_DIR)
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 main.go
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 main.go
	# macOS AMD64 (Intel)
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 main.go
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe main.go
	@echo "Cross-compilation complete in $(BUILD_DIR)/"
